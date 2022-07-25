# frozen_string_literal: true

# Ratelimiter is responsible for encapsulating the data for the ratelimit.go file.
class Ratelimiter
  attr_reader \
    :endpoint

  def initialize(endpoint)
    @endpoint = endpoint
  end
end
