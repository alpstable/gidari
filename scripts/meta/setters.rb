module Setters
  private

  def get_body_setter(field)
    "internal.HTTPBodyFragment(body,  \"#{field.identifier}\", #{OPTIONS_ALIAS}.#{field.go_variable_name})"
  end

  def get_query_param_setter(field)
    sig = {
      '[]string' => "\ninternal.HTTPQueryEncodeStrings",

      'bool' => 'internal.HTTPQueryEncodeBool',
      'int32' => 'internal.HTTPQueryEncodeInt32',
      'int' => 'internal.HTTPQueryEncodeInt',
      'uint8' => 'internal.HTTPQueryEncodeUint8',
      'string' => 'internal.HTTPQueryEncodeString',
      'time.Time' => 'internal.HTTPQueryEncodeTime'
    }[field.go_type]

    sig = 'internal.HTTPQueryEncodeStringer' if sig.nil?

    adr = field.required && !field.go_slice? ? '&' : ''
    sig + "(req, \"#{field.identifier}\", #{adr}#{OPTIONS_ALIAS}.#{field.go_field_name})"
  end

  public

  def option_setters(endpoint)
    return unless endpoint.params?

    endpoint.all_params.dup.map do |field|
      variable_name = field.go_variable_name
      variable_name = 't' if variable_name == 'type'
      sig_go_type = field.ptr_go_type.include?('[]*') ? "[]*#{field.go_type.dup.gsub('[]', '')}" : field.go_type
      struct = "#{endpoint.go_model_name}Options"

      r_ptr = "func (#{OPTIONS_ALIAS} *#{struct})"
      r_name = "Set#{field.go_protofield_name}"
      r_sig = "#{r_name}(#{variable_name} #{sig_go_type})"
      r_ret = "*#{struct}"
      comment = "#{r_name} sets the #{field.go_protofield_name} field on #{struct}."
      comment += "  #{field.description}" unless field.description.nil?

      logic = "#{OPTIONS_ALIAS}.#{field.go_protofield_name} = #{field.ptr_go_variable};"
      logic += "return #{OPTIONS_ALIAS}"

      { setter: "\n#{format_go_comment(comment)}\n#{r_ptr} #{r_sig} #{r_ret} {\n#{logic}\n}\n", name: struct }
    end
  end

  def option_body_setter(endpoint)
    return unless endpoint.params?

    name = 'EncodeBody'
    struct = "#{endpoint.go_model_name}Options"
    sig = "func (#{OPTIONS_ALIAS} *#{struct}) #{name}() (buf io.Reader, err error)"
    top = false

    if endpoint.body?
      setters = endpoint.body.dup.map { |field| get_body_setter(field) }.compact.flatten.join(';')
      logic = "{\n"
      logic += "if #{OPTIONS_ALIAS} != nil {;"
      logic += 'body := make(map[string]interface{});'
      logic += setters
      logic += ';raw, err := json.Marshal(body);'
      logic += 'if err == nil {; buf = bytes.NewBuffer(raw); };'
      logic += '};'
      logic += "return;\n}\n"
    else
      top = true
      logic = '{return}'
    end

    { setter: "\n#{sig} #{logic}", name: struct, top: top }
  end

  def option_query_params_setter(endpoint)
    return unless endpoint.params?

    name = 'EncodeQuery'
    struct = "#{endpoint.go_model_name}Options"
    sig = "func (#{OPTIONS_ALIAS} *#{struct}) #{name}(req *http.Request)"
    top = false

    if endpoint.query_params?
      setters = endpoint.query_params.dup.map { |field| get_query_param_setter(field) }.compact.flatten
      logic = "{\n"
      logic += "if #{OPTIONS_ALIAS} != nil { #{setters.join("\n")} };"
      logic += "return;\n}\n"
    else
      top = true
      logic = '{return}'
    end

    { setter: "\n#{sig} #{logic}", name: struct, top: top }
  end
end
