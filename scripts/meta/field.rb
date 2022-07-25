# frozen_string_literal: true

require_relative 'comment'
require_relative 'inflector'

# Field holds state concerning endpoints given by the meta/schema json files.  It encapsulates methods for manupilating
# this data for various use cases in go, such as structs, functions, closures, etc.
class Field
  attr_reader \
    :datetime_layout,
    :deserializer,
    :identifier,
    :go_type,
    :go_field_name,
    :go_field_tag,
    :description,
    :hash,
    :required,
    :go_comment,
    :sql_identifier

  include Comment
  include Inflector

  GO_TYPES = %w(
    string
    bool
    time.Time
    int
    []string
    float64
  ).freeze

  def initialize(hash)
    return if hash.nil?

    @hash = hash
    @datetime_layout = hash[:datetimeLayout] || 'time.RFC3339Nano'
    @deserializer = hash[:unmarshaler]
    @identifier = hash[:identifier]
    @sql_identifier = inflector.underscore(hash[:identifier])
    @go_type = hash[:goType]
    @go_field_name = inflector.camelize_upper(hash[:identifier].dup.gsub('.', '_').gsub('-', '_'))
    @go_field_tag = inflector.camelize_lower("#{hash[:identifier]}_json_tag")
    @description = hash[:description] || ''
    @go_comment = format_go_comment(@description) unless @description.nil? || @description == ''
    @required = hash[:required]
  end

  def custom_type?
    return true if !datetime_layout.nil? && datetime_layout != 'time.RFC3339Nano'
    return true unless deserializer.nil?

    false
  end

  def go_slice?
    @go_type.include?('[]')
  end

  def go_struct?
    return false if GO_TYPES.include?(hash[:goType])
    return false if go_type.include?('scalar')

    true
  end

  def go_protofield_name
    return go_field_name unless go_struct?

    go_field_name.to_s
  end

  def go_variable_name
    name = inflector.camelize_lower(@go_field_name)

    # `type` is a go keyword, the convention will be to replace it with `typ`.
    return 'typ' if name == 'type'

    name
  end

  def ptr_go_type
    return go_type if required

    if go_type.include?('[]')
      "[]*#{go_type.dup.gsub('[]', '')}"
    else
      "*#{go_type}"
    end
  end

  def ptr_go_variable
    return go_variable_name if required || go_type.include?('[]')

    "&#{go_variable_name}"
  end
end
