require 'minitest/autorun'

Dir.glob('test_*.rb') { |f| require_relative(f) }
