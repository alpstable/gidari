# frozen_string_literal: true

require_relative 'comment'
require_relative 'field'
require_relative 'inflector'

# Enum is a data fromatter and accessor for enum types. This class is used to encapsulate data for building the types.go
# file in web API packages.
class Enum
  attr_reader \
    :description,
    :go_comment,
    :go_type_name,
    :go_var_name,
    :identifier,
    :pluralize,
    :pluralize_var,
    :values

  include Comment
  include Inflector

  def initialize(hash = {})
    @description		= hash[:description]
    @go_comment	= format_go_comment(@description) unless @description.nil?
    @identifier 		= hash[:identifier]

    parse_go_plural(hash)
    parse_go_singular(hash)
    parse_values(hash)
  end

  private

  def parse_go_plural(hash)
    return if hash[:pluralize].nil?

    @pluralize	= inflector.camelize_upper(hash[:pluralize].dup.gsub('.', '_').gsub('-', '_'))
    @pluralize_var	= inflector.camelize_lower(pluralize)
  end

  def parse_go_singular(hash)
    @go_type_name 	= inflector.camelize_upper(hash[:identifier].dup.gsub('.', '_').gsub('-', '_'))
    @go_var_name 		= inflector.camelize_lower(go_type_name)
  end

  def parse_values(hash)
    @values = hash[:values].map do |val|
      Field.new(val)
    end
  end
end
