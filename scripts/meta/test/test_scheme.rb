require 'minitest/autorun'
require_relative '../scheme'

describe Scheme do
  before do
    @filename = "#{File.dirname(__FILE__)}/model/test_model_1.json"
    @scheme = Scheme.new(@filename)
  end

  describe 'when reading the filename' do
    it 'must respond with initialized filename' do
      _(@scheme.filename).must_equal(@filename)
    end
  end

  describe 'when parsing json from the filename' do
    it 'must respond with a hash of data' do
      _(@scheme.api).must_equal('test_api')
      _(@scheme.go_model_filename).must_equal('test_model.go')
      _(@scheme.go_model_name).must_equal('TestModel')
    end
  end
end
