# frozen_string_literal: true

require_relative 'const'

# ModelUnmarshaler is responsible for the logic to write the unmarshaler methods
# for models that require specific-use-case non-native go-type support.
module ModelUnmarshaler
  private

  def constantize_json_go_tags
    c = fields.dup.map { |field| "#{field.go_field_tag} = \"#{field.identifier}\"" }
    "const(#{c.join(';')})"
  end

  def struct_access_variable(field)
    "#{go_model_variable_name}.#{field.go_protofield_name}"
  end

  def unmarshal_fn_signature(field)
    "#{field.go_field_tag}, &#{struct_access_variable(field)}"
  end

  def custom_deserializer(field)
    return nil unless !field.deserializer.nil? && field.deserializer == 'UnmarshalFloatString'

    "\ndata.UnmarshalFloatString(#{unmarshal_fn_signature(field)})"
  end

  def time_deserializer(field)
    sig = [field.datetime_layout, field.go_field_tag, "&#{struct_access_variable(field)}"]
    ["err = data.UnmarshalTime(#{sig.join(',')})", Const::RETURN_ERR].join(';')
  end

  def scalar_deserializer(field, sig)
    "\ndata.Unmarshal#{field.go_type.dup.gsub('scalar.', '')}(#{sig})"
  end

  def slice_deserializer(field)
    accessor = struct_access_variable(field)
    data_accessor = "data.Value(#{field.go_field_tag})"
    type = field.go_type.dup.gsub!('[]', '').gsub!('*', '')
    marshal_logic = "bytes, _ := json.Marshal(item); obj := #{type}{};"
    unmarshal_logic = 'if err := json.Unmarshal(bytes, &obj); err != nil {return err};'
    logic = "#{marshal_logic} #{unmarshal_logic} #{accessor} = append(#{accessor}, &obj)"
    loop_logic = "for _, item := range #{data_accessor}.([]interface{}) {#{logic}}"

    "\nif v := #{data_accessor}; v != nil {#{loop_logic}}"
  end

  def struct_deserializer(field, sig)
    accessor = struct_access_variable(field)
    init = "#{accessor} = #{field.go_type}{}"

    "\n#{init}; if err := data.UnmarshalStruct(#{sig}); err != nil {return err}"
  end

  def get_deserializer(field, sig)
    # if the deserializer is passed into the field via the schema, then just type it exactly.  Otherwise, just use the
    # default type deserializer.
    unless field.deserializer.nil?
      sig = [field.go_field_tag, "&#{struct_access_variable(field)}"]
      return ["err = data.#{field.deserializer}(#{sig.join(',')})", Const::RETURN_ERR].join(';')
    end
    {
      'string' => "\ndata.UnmarshalString(#{sig})",
      'bool' => "\ndata.UnmarshalBool(#{sig})",
      'time.Time' => "\n#{time_deserializer(field)}",
      'int' => "\ndata.UnmarshalInt(#{sig})",
      'int32' => "\ndata.UnmarshalInt32(#{sig})",
      '[]string' => "\ndata.UnmarshalStringSlice(#{sig})",
      'float64' => "\ndata.UnmarshalFloat(#{sig})"
    }[field.go_type]
  end

  def generic_deserializer(field)
    sig = unmarshal_fn_signature(field)
    deserializer = get_deserializer(field, sig)
    return deserializer unless deserializer.nil?
    return scalar_deserializer(field, sig) if field.go_type.include?('scalar')
    return slice_deserializer(field) if field.go_type.include?('[]')

    struct_deserializer(field, sig)
  end

  def deserializers
    fields.dup.map do |field|
      custom_deserializer(field) || generic_deserializer(field)
    end.sort.join(';')
  end

  public

  def unmarshaler
    comment = format_go_comment("UnmarshalJSON will deserialize bytes into a #{go_model_name} model")
    return "\n" + [comment, "\n" + custom_unmarshaler].join('') unless custom_unmarshaler.nil?
    return '' unless non_struct.nil?
    return '' if fields.empty?
    return '' unless custom_type?

    serial = "\ndata, err := serial.NewJSONTransform(d); #{Const::RETURN_ERR}"
    fn = [constantize_json_go_tags, serial, deserializers]

    comment = format_go_comment("UnmarshalJSON will deserialize bytes into a #{go_model_name} model")

    reciever = "\nfunc (#{go_model_variable_name} *#{go_model_name}) "
    signature = "UnmarshalJSON(d []byte) error {#{fn.join('')}; return nil}"
    function = reciever + signature

    "\n" + [comment, function].join('')
  end
end
