# frozen_string_literal: true

# GoClient will build the HTTP client methods for a schema.
module GoClient
  def self.apis(schema)
    tree = (proc { Hash.new { |hash, key| hash[key] = [] } }).call
    schema.each { |scheme| tree[scheme.api] << scheme }
    tree
  end

  def self.receivers(schema)
    schema.dup.map(&:client_receivers).flatten.compact.sort_by { |r| r[:name] }.map { |r| r[:receiver] }.join("\n")
  end

  def self.write(schema)
    apis(schema).each do |api, api_schema|
      path = Pathname.new(PARENT_DIR).join(api)
      Dir.chdir(path.to_s) do
        File.open(CLIENT_FILENAME, 'w') do |f|
          f.write("package #{api}")
          f.write("\nimport \"github.com/alpine-hodler/driver/web/scalar\";")
          f.write("\nimport \"github.com/alpine-hodler/driver/internal/serial\";")
					f.write("\nimport \"github.com/alpine-hodler/driver/internal/client\";")
					f.write("\nimport \"github.com/alpine-hodler/driver/internal\";")
					f.write("\nimport \"golang.org/x/time/rate\";")
          f.write(GEN_MSG)
          f.write(receivers(api_schema))
        end
        `/go/bin/goimports -w #{CLIENT_FILENAME}`
      end
    end
  end
end
