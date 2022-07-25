# frozen_string_literal: true

# PostAuthority holds endpoint data by api
class Path
  attr_reader \
    :endpoint

  def initialize(endpoint)
    @endpoint = endpoint
  end
end
