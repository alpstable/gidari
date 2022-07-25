# frozen_string_literal: true

require 'json'
require 'json-schema'

require_relative 'comment'
require_relative 'field'
require_relative 'endpoint'
require_relative 'model_unmarshaler'
require_relative 'go_struct'
require_relative 'option'
require_relative 'setters'
require_relative 'go_http'
require_relative 'inflector'

# Scheme is the class encapsulation of a single json file in the meta/schema
# directory
class Scheme
  attr_reader \
    :api,
    :description,
    :go_comment,
    :go_model_filename,
    :go_model_name,
    :model,
    :model_only,
    :filename,
    :fields,
    :endpoints,
    :go_model_variable_name,
    :non_struct,
    :package,
    :custom_unmarshaler

  include Comment
  include ModelUnmarshaler
  include Option
  include GoStruct
  include Setters
  include GoHTTP
  include Inflector

  def initialize(filename)
    file = File.read(filename)
    hash = JSON.parse(file, symbolize_names: true)
    validate(hash, filename)

    @api = hash[:api]
    @package = "pacakge #{api}"
    @description = hash[:modelDescription]
    @filename = filename
    @model = hash[:model].to_s
    @model_only = hash[:modelOnly] || false
    @non_struct = hash[:nonStruct]
    @go_comment = format_go_comment(@description)
    @go_model_filename = "#{@model}.go"
    @go_model_name = inflector.camelize_upper(@model)
    @go_model_variable_name = inflector.camelize_lower(@go_model_name)
    @custom_unmarshaler = hash[:customUnmarshaler]

    @fields = hash[:modelFields].map { |field| Field.new(field) }
    @endpoints = (hash[:endpoints] || []).map { |ep| Endpoint.new(api, ep) }
  end

  def custom_type?
    fields.any?(&:custom_type?)
  end

  def validate(hash, filename)
    schema = JSON.parse(File.read("#{File.dirname(__FILE__)}/schema/schema.json"))
    e = JSON::Validator.fully_validate(schema, hash)
    raise "Schema Error for #{filename}: #{e}" unless e.empty?
  end
end
