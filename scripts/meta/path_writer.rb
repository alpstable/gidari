# frozen_string_literal: true

require_relative 'comment'
require 'string_inflection'
using StringInflection

# PathWriter is responsible for methods that generate the endpoint.go code in
# web packages
module PathWriter
  MSG = "\n// * This is a generated file, do not edit\n"

  def self.endpoint_consts(endpoints)
    consts = ["_ #{POST_AUTHORITY_TYPE_ALIAS} = iota"] | endpoints.dup.map { |ep| ep.go_const }.sort
    "const(#{consts.join(';')})"
  end

  def self.get_function(endpoints)
    mappings = endpoints.dup.map { |ep| "#{ep.go_const}: get#{ep.enum_root.to_pascal}Path," }
    comment = Comment.u_format_go_comment("Get takes an #{POST_AUTHORITY_TYPE_ALIAS} const and #{POST_AUTHORITY_TYPE_ALIAS} arguments to parse the URL #{POST_AUTHORITY_TYPE_ALIAS} path.")
    rec = "\n#{comment}\nfunc (#{POST_AUTHORITY_ALIAS} #{POST_AUTHORITY_TYPE_ALIAS})"
    sig = "Path(#{URI_BUILDER_ALIAS} map[string]string) string"

    map = "map[#{POST_AUTHORITY_TYPE_ALIAS}]func(map[string]string) string"
    wrapper = "#{map}{\n#{mappings.join("\n")}\n}"

    "#{rec} #{sig} {return #{wrapper}[#{POST_AUTHORITY_ALIAS}](#{URI_BUILDER_ALIAS})};\n"
  end

  def self.scope_function(endpoints)
    mappings = endpoints.dup.map { |ep| ep.scope.nil? ? nil : "#{ep.go_const}: \"#{ep.scope}\"," }.flatten.compact
    # comment = Comment.u_format_go_comment("Get takes an #{POST_AUTHORITY_TYPE_ALIAS} const and #{POST_AUTHORITY_TYPE_ALIAS} arguments to parse the URL #{POST_AUTHORITY_TYPE_ALIAS} path.")
    rec = "\nfunc (#{POST_AUTHORITY_ALIAS} #{POST_AUTHORITY_TYPE_ALIAS})"
    sig = 'Scope() string'

    map = "map[#{POST_AUTHORITY_TYPE_ALIAS}]string"
    wrapper = "#{map}{\n#{mappings.join("\n")}\n}"

    "#{rec} #{sig} {return #{wrapper}[#{POST_AUTHORITY_ALIAS}]};\n"
  end

  def self.pkg_name(api)
    "package #{api}"
  end

  def self.join_paths(endpoint)
    "path.Join(#{endpoint.path_parts.dup.map { |pp| pp.go_arg }.join(',')})"
  end

  def self.path_functions(endpoints)
    endpoints.dup.map do |endpoint|
      var = '_'
      var = URI_BUILDER_ALIAS if endpoint.path_parts? || endpoint.query_params?

      sig = "func get#{endpoint.enum_root.to_pascal}Path(#{var} map[string]string) string"
      logic = "{\nreturn #{join_paths(endpoint)}}"
      desc = Comment.u_format_go_comment(endpoint.description)
      "\n#{desc}\n #{sig} #{logic}"
    end.join("\n\n")
  end

  def self.apis(paths)
    tree = (proc { Hash.new { |hash, key| hash[key] = [] } }).call
    paths.each do |path|
      next if path.endpoint.nil?

      tree[path.endpoint.api] << path.endpoint
    end
    tree
  end

  def self.write(paths)
    apis(paths).each do |api, endpoints|
      endpoints = endpoints.sort_by { |ep| ep.enum_root }
      Dir.chdir(PARENT_DIR + "/#{api}") do
        File.open(POST_AUTHORITY_FILENAME, 'w') do |f|
          f.write(pkg_name(api))
          f.write("\nimport \"github.com/alpine-hodler/driver/internal/client\";")
          f.write("\nimport \"github.com/alpine-hodler/driver/internal\";")
          f.write(MSG)
          f.write("\ntype #{POST_AUTHORITY_TYPE_ALIAS}  uint8;")
          f.write(endpoint_consts(endpoints))
          f.write(path_functions(endpoints))
          f.write(get_function(endpoints))
          f.write(scope_function(endpoints))
        end

        # In addition to fixing imports, goimports also formats your code in the
        # same style as gofmt so it can be used as a replacement for your editor's
        # gofmt-on-save hook.
        `/go/bin/goimports -w #{POST_AUTHORITY_FILENAME}`
      end
    end
  end
end
