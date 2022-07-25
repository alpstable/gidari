#!/usr/local/bin/ruby -w
# frozen_string_literal: true

require_relative 'scheme'
require_relative 'model_writer'
require_relative 'path'
require_relative 'path_writer'
require_relative 'option'
require_relative 'go_client'
require_relative 'path_part'
require_relative 'types'
require_relative 'types_writer'
require_relative 'ratelimiter'
require_relative 'ratelimiter_writer'

URI_BUILDER_ALIAS = 'params'
CLIENT_ALIAS = 'c'
CLIENT_FILENAME = 'http.go'
CLIENT_PKG = 'client'
CLIENT_STRUCT_NAME = 'Client'
GEN_MSG = "\n// * This is a generated file, do not edit\n"
OPTIONS_ALIAS = 'opts'
OPTIONS_FILENAME = 'options.go'
PARENT_DIR = File.expand_path('..', Dir.pwd)
POST_AUTHORITY_ALIAS = 'p'
POST_AUTHORITY_FILENAME = 'path.go'
POST_AUTHORITY_TYPE_ALIAS = 'rawPath'
RETURN_ALIAS = 'm'
TOOLS_PKG = 'tools'

def generate_models
  paths = []
  schema = []
  types = []
  ratelimiters = []
  Dir.glob("#{File.dirname(__FILE__)}/schema/*").each do |dir|
    Dir.glob("#{dir}/*.json").each do |filename|
      if filename.include?('/types.json')
        types << Types.new(filename)
      else
        scheme = Scheme.new(filename)
        schema << scheme
        scheme.endpoints.each do |ep|
          ratelimiters << Ratelimiter.new(ep)
          paths << Path.new(ep)
        end
      end
    end
  end

  ModelWriter.write(schema)
  Option.write(schema)
  GoClient.write(schema)
  TypesWriter.write(types)
  RatelimitWriter.write(ratelimiters)
  PathWriter.write(paths)
end

generate_models
