require 'minitest/autorun'
require_relative '../endpoint'

URI_BUILDER_ALIAS = 'builder'
TOOLS_PKG = 'tools'

describe Endpoint do
  before do
    @endpoint1 = Endpoint.new('test_api', {
                                path: 'some/{path_id}',
                                enumRoot: 'someEnumRoot',
                                description: 'some description',
                                queryParams: [
                                  {
                                    identifier: 'query_param1',
                                    goType: 'string',
                                    required: true
                                  },
                                  {
                                    identifier: 'query_param2',
                                    goType: 'string'
                                  }
                                ]
                              })

    @endpoint2 = Endpoint.new('test_api', {
                                path: 'some/{path_id}',
                                enumRoot: 'someEnumRoot',
                                description: 'some description'
                              })
  end

  describe 'when parsing data from the list' do
    it 'must initialize accessors' do
      _(@endpoint1.path).must_equal('some/{path_id}')
      _(@endpoint1.enum_root).must_equal('someEnumRoot')
      _(@endpoint1.go_const).must_equal('someEnumRootPostAuthority')
      _(@endpoint1.description).must_equal('some description')

      _(@endpoint1.path_parts.length).must_equal(2)
      _(@endpoint1.path_parts[0].go_arg).must_equal('"/some"')
      _(@endpoint1.path_parts[1].go_arg).must_equal("builder.Get(internal.URIBuilderComponentPath, \"path_id\")")

      _(@endpoint1.query_params.length).must_equal(2)
      _(@endpoint2.query_params.length).must_equal(0)
    end
  end
end
