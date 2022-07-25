# frozen_string_literal: true

require 'json'
require 'json-schema'

require_relative 'enum'

# Types is responsible for building the types go file for an API.
class Types
  attr_reader \
    :api,
    :enums

  def initialize(filename)
    file = File.read(filename)
    hash = JSON.parse(file, symbolize_names: true)
    # validate(hash, filename)

    @api = hash[:api]
    @enums = hash[:enums].dup.each.map { |enum| Enum.new(enum) }
  end
end
