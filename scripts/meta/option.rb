# frozen_string_literal: true

# Model will generate the go files with structs that inherit the protomodel
# fields and methods, but is open to modification.
module Option
  PKG = 'option'
  PKG_DEC = "package #{PKG}"
  MSG = "\n// * This is a generated file, do not edit\n"

  def self.apis(schema)
    tree = (proc { Hash.new { |hash, key| hash[key] = [] } }).call
    schema.each { |scheme| tree[scheme.api] << scheme }
    tree
  end

  def self.structs(schema)
    schema.dup.map do |scheme|
      scheme.endpoints.dup.map { |ep| scheme.options_struct(ep) }
    end.join('')
  end

  def self.setters(schema)
    schema.dup.map do |scheme|
      scheme.endpoints.dup.map { |ep| scheme.option_setters(ep) }
    end.flatten.compact.sort_by { |r| r[:name] }.map { |r| r[:setter] }.join("\n")
  end

  def self.body_setters(schema)
    schema.dup.map do |scheme|
      scheme.endpoints.dup.map { |ep| scheme.option_body_setter(ep) }
    end.flatten.compact.delete_if { |h| h[:top] == true }.sort_by { |r| r[:name] }.map { |r| r[:setter] }.join("\n")
  end

  def self.body_setters_top(schema)
    schema.dup.map do |scheme|
      scheme.endpoints.dup.map { |ep| scheme.option_body_setter(ep) }
    end.flatten.compact.delete_if { |h| h[:top] == false }.sort_by { |r| r[:name] }.map { |r| r[:setter] }.join(";")
  end

  def self.query_param_setters(schema)
    schema.dup.map do |scheme|
      scheme.endpoints.dup.map { |ep| scheme.option_query_params_setter(ep) }
    end.flatten.compact.delete_if { |h| h[:top] == true }.sort_by { |r| r[:name] }.map { |r| r[:setter] }.join("\n")
  end

  def self.query_param_setters_top(schema)
    schema.dup.map do |scheme|
      scheme.endpoints.dup.map { |ep| scheme.option_query_params_setter(ep) }
    end.flatten.compact.delete_if { |h| h[:top] == false }.sort_by { |r| r[:name] }.map { |r| r[:setter] }.join(";")
  end

  def self.write(schema)
    apis(schema).each do |api, api_schema|
      path = Pathname.new(PARENT_DIR).join(api)
      Dir.chdir(path.to_s) do
        File.open('options.go', 'w') do |f|
          f.write("package #{api}")
          f.write("\nimport \"github.com/alpine-hodler/driver/web/scalar\";")
          f.write("\nimport \"github.com/alpine-hodler/driver/internal/serial\";")
          f.write("\nimport \"time\";")
          f.write("\nimport \"github.com/alpine-hodler/driver/internal\";")
          f.write(Option::MSG)
          f.write(structs(api_schema))
          f.write(body_setters_top(api_schema))
          f.write(query_param_setters_top(api_schema))
          f.write(setters(api_schema))
          f.write(body_setters(api_schema))
          f.write(query_param_setters(api_schema))
        end
        `/go/bin/goimports -w #{OPTIONS_FILENAME}`
      end
    end
  end
end
