# frozen_string_literal: true

require_relative 'path_part'
require_relative 'field'
require_relative 'inflector'

# Endpoint holds state concerning endpoints given by the mete/schema json files
class Endpoint
  attr_reader \
    :hash,
    :path,
    :enum_root,
    :go_const,
    :path_parts,
    :description,
    :query_params,
    :go_model_name,
    :go_query_param_filename,
    :slice,
    :body,
    :return_type,
    :http_method,
    :documentation,
		:scope,
		:rate_limit,
		:api,
		:go_rl_const

  include Inflector

  def initialize(api, hash)
    return if hash.nil?

    @hash = hash
    @path = hash[:path]
    @enum_root = inflector.camelize_lower(hash[:enumRoot])
    @go_const = "#{inflector.camelize_lower(enum_root)}Path"
		@go_rl_const = "#{inflector.camelize_lower(enum_root)}Ratelimiter"
    @description = hash[:description] || ''
    @slice = hash[:slice]
    @return_type = hash[:returnType]
    @http_method = hash[:httpMethod]
    @documentation = hash[:documentation]
    @go_model_name = inflector.camelize_upper(enum_root)
    @go_query_param_filename = "#{api}_#{hash[:enum_roof]}.go"
		@scope = hash[:scope]
		@rate_limit = hash[:rateLimit].to_i
		@api = api

    set_path_parts
    set_query_params
    set_body
  end

  def all_params
    @query_params + @body
  end

  def body?
    !@body.empty?
  end

  def path_parts?
    !@path_parts.empty?
  end

  def params?
    body? || query_params?
  end

  def path_params
    return [] if path_parts.empty?

    path_parts.map { |part| part.path_param? ? part : nil }.compact
  end

  def query_params?
    !@query_params.empty?
  end

  def no_assignment?
    @return_type == 'none'
  end

  private

  def set_path_parts
    @path_parts = []
    first_part_set = false
    path.split('/').each do |part|
      next if part == ''

      part = '/' + part unless first_part_set
      first_part_set = true
      @path_parts << PathPart.new(part)
    end
  end

  def set_body
    @body = []
    (hash[:body] || []).each do |subhash|
      @body << Field.new(subhash)
    end
  end

  def set_query_params
    @query_params = []
    (hash[:queryParams] || []).each do |subhash|
      @query_params << Field.new(subhash)
    end
  end
end
