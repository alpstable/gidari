require 'string_inflection'
using StringInflection

class PathPart
  attr_reader \
    :name,
    :go_var,
    :go_arg,
    :go_type

  GO_TYPES = %w(
    string
    bool
    time.Time
    int
    []string
    float64
  ).freeze

  def initialize(name)
    if name.include?(':')
      @name = name.split(':')[0]
      @go_type = name.split(':')[1].dup.gsub('}', '').gsub('{', '')
    else
      @name = name
      @go_type = 'string'
    end
    @go_var = self.name.to_camel

    set_go_arg
  end

  def param_go_var force_string = false
    return '' unless path_param?

    n = name.dup.gsub('}', '').gsub('{', '').to_camel
    return "#{n}.String()" if !@go_type.nil? && !GO_TYPES.include?(@go_type) && force_string

    n
  end

  def param_name
    name.dup.gsub('}', '').gsub('{', '')
  end

  def path_param?
    return true if name.include?('{') || name.include?('}')

    false
  end

  private

  def set_go_arg
    @go_arg = if path_param?
                "params[\"#{name.dup.gsub!('{', '').gsub('}', '')}\"]"
              else
                "\"#{name}\""
              end
  end
end
