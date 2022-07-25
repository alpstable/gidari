# frozen_string_literal: true

require 'string_inflection'
using StringInflection

# HTTP will build the HTTP client methods for a scheme.
module GoHTTP
  private

  def path_param_args(endpoint)
    return nil unless endpoint.path_parts?

    endpoint.path_params.dup.map { |part| "#{part.param_go_var} #{part.go_type}" }
  end

  def option_arg(endpoint)
    return nil unless endpoint.params?

    "opts *#{endpoint.go_model_name}Options"
  end

  def sig_args(endpoint)
    args = [path_param_args(endpoint), option_arg(endpoint)].flatten.compact
    return "(#{args[0]})" if args.length == 1

    "(#{args.join(',')})"
  end

  def return_args(endpoint)
    return 'error' if endpoint.no_assignment?

    val = if endpoint.return_type.nil?
            endpoint.slice ? "[]*#{go_model_name}" : "*#{go_model_name}"
          else
            endpoint.return_type
          end

    "(#{RETURN_ALIAS} #{val}, _ error)"
  end

  def declare_request(endpoint)
    opts = endpoint.params? ? 'opts' : 'nil'
    "req, _ := internal.HTTPNewRequest(\"#{endpoint.http_method}\", \"\", #{opts});"
  end

  def params_function(endpoint)
    map_input = endpoint.path_parts.dup.map do |path_part|
      next unless path_part.path_param?

      "\"#{path_part.param_name}\": #{path_part.param_go_var(true)},"
    end
    return 'nil' if map_input.compact.empty?

    "map[string]string{\n#{map_input.flatten.compact.join("\n")}}"
  end

  def return_stmt(endpoint)
    opts_var = endpoint.query_params? ? OPTIONS_ALIAS : 'nil'

    if endpoint.no_assignment?
      'return internal.HTTPFetch(nil,' +
        "internal.HTTPWithClient(#{CLIENT_ALIAS}.Client),\n" +
        "internal.HTTPWithEncoder(#{opts_var}),\n" +
        "internal.HTTPWithEndpoint(#{endpoint.go_const}),\n" +
        "internal.HTTPWithParams(#{params_function(endpoint)}),\n" +
        "internal.HTTPWithRatelimiter(getRateLimiter(#{endpoint.go_rl_const}, #{endpoint.rate_limit})),\n" +
        'internal.HTTPWithRequest(req))'
    else
      "return #{RETURN_ALIAS}, internal.HTTPFetch(&#{RETURN_ALIAS}," +
        "internal.HTTPWithClient(#{CLIENT_ALIAS}.Client),\n" +
        "internal.HTTPWithEncoder(#{opts_var}),\n" +
        "internal.HTTPWithEndpoint(#{endpoint.go_const}),\n" +
        "internal.HTTPWithParams(#{params_function(endpoint)}),\n" +
        "internal.HTTPWithRatelimiter(getRateLimiter(#{endpoint.go_rl_const}, #{endpoint.rate_limit})),\n" +
        'internal.HTTPWithRequest(req))'
    end
  end

  public

  # client_receivers will generate all of the Go `cleint` struct receivers from a scheme's endpoints.
  def client_receivers
    return nil if endpoints.empty?

    receiver_alias = "(#{CLIENT_ALIAS} *#{CLIENT_STRUCT_NAME})"
    endpoints.dup.map do |endpoint|
      r = "\n" + format_go_comment_with_source(endpoint.description, endpoint.documentation)
      r += "\nfunc"
      r += receiver_alias
      r += ' ' + endpoint.enum_root.to_pascal
      r += sig_args(endpoint)
      r += ' ' + return_args(endpoint)
      r += ' ' + "{\n #{declare_request(endpoint)};#{return_stmt(endpoint)} \n}"

      { receiver: r, name: endpoint.enum_root }
    end
  end
end
