# frozen_string_literal: true

require_relative 'const'

# ModelWriter will generate the go models for unmarshaling web API data into.
module ModelWriter
  def self.apis(schema)
    tree = (proc { Hash.new { |hash, key| hash[key] = [] } }).call
    schema.each { |scheme| tree[scheme.api] << scheme }
    tree
  end

  def self.structs(schema)
    schema.dup.map(&:model_struct).join('')
  end

  def self.unmarshallers(schema)
    schema.dup.map(&:unmarshaler).join("\n")
  end

  def self.write_file(file_, api, schema)
    file_.write(Const.package(api))
    file_.write(Const.import(Const::SERIAL_PKG))
    file_.write(Const::GEN_MSG)
    file_.write(structs(schema))
    file_.write(unmarshallers(schema))
  end

  def self.write(schema)
    apis(schema).each do |api, api_schema|
      path = Pathname.new(Const::PARENT_DIR).join(api)
      Dir.chdir(path.to_s) do
        File.open(Const::MODEL_FILENAME, 'w') do |f|
          write_file(f, api, api_schema)
        end
        Const.go_fmt(Const::MODEL_FILENAME)
      end
    end
  end
end
