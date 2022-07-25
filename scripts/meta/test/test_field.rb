require 'minitest/autorun'
require_relative '../field'

describe Field do
  before do
    @field1 = Field.new({
                          identifier: 'id',
                          goType: 'string',
                          unmarshaler: 'UnmarshalFloatString'
                        })

    @field2 = Field.new({
                          identifier: 'some_datetime',
                          goType: 'SomeFancyStruct',
                          datetimeLayout: 'iso8601',
                          description: 'some description'
                        })
  end

  describe 'when parsing data from the list' do
    it 'must initialize accessors' do
      _(@field1.identifier).must_equal('id')
      _(@field1.go_type).must_equal('string')
      _(@field1.datetime_layout).must_equal('time.RFC3339Nano')
      _(@field1.deserializer).must_equal('UnmarshalFloatString')
      _(@field1.go_field_name).must_equal('ID')
      _(@field1.go_field_tag).must_equal('IDJSONTag')
      _(@field1.description).must_equal('')

      _(@field2.datetime_layout).must_equal('iso8601')
      _(@field2.go_field_tag).must_equal('someDatetimeJSONTag')
      _(@field2.description).must_equal('some description')
    end
  end

  describe 'when determining if a hash represents a go struct' do
    it 'must return false when it is not a go struct' do
      _(@field1.go_struct?).must_equal(false)
    end

    it 'must return true when it is a go struct' do
      _(@field2.go_struct?).must_equal(true)
    end
  end

  describe 'when a field represents a go struct' do
    it 'must return a protofield name' do
      _(@field2.go_protofield_name).must_equal('SomeDatetime')
    end
  end
end
