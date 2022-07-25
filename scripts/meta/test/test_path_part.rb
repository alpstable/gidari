require 'minitest/autorun'
require_relative '../path_part'

describe PathPart do
  # before do
  #   @path_part1 = PathPart.new('test_path_part')
  #   @path_part2 = PathPart.new('{path_param}')
  # end

  # describe 'when parsing data from the list' do
  #   it 'must initialize accessors' do
  #     _(@path_part1.name).must_equal('test_path_part')
  #     _(@path_part1.go_var).must_equal('testPathPart')
  #     _(@path_part1.go_arg).must_equal('"test_path_part"')

  #     _(@path_part2.go_arg).must_equal("builder.Get(internal.URIBuilderComponentPath, \"path_param\")")
  #   end
  # end

  # describe 'when using path_param?' do
  #   it 'must return false when the value is not a path param' do
  #     _(@path_part1.path_param?).must_equal(false)
  #   end

  #   it 'must return true when the value is a path param' do
  #     _(@path_part2.path_param?).must_equal(true)
  #   end
  # end
end
