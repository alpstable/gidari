# frozen_string_literal: true

require 'dry/inflector'

# Inflector returns a dry inflector that can be assigned acronyms to inflect when converting a string between different
# cases such as camel, snake, and pascal.
module Inflector
  def inflector
    Dry::Inflector.new do |inflections|
      inflections.acronym 'ID', 'FIAT', 'SEPA', 'UK', 'SWIFT', 'ACH', 'STP', 'URL', 'GTC', 'IOC', 'FOK', 'GTT', 'STP'
    end
  end
end
