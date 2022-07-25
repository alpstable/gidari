# frozen_string_literal: true

# GoStruct will build the go struct scheme.
module GoStruct
  def model_struct
    return "\n#{go_comment}\ntype #{go_model_name} #{non_struct}\n" unless non_struct.nil?

    literals = []
    fields.dup.sort_by { |f| f.go_field_name }.each do |field|
      literal = field.description.empty? ? '' : "\n#{format_go_comment(field.description)}\n"
      literal += "#{field.go_protofield_name} #{field.go_type}"
      literal += "`bson:\"#{field.sql_identifier}\" json:\"#{field.identifier}\" sql:\"#{field.sql_identifier}\"`"
      literals << literal
    end
    literals = literals.sort.join("\n")
    "\n#{go_comment}\ntype #{go_model_name} struct {#{literals}}\n"
  end

  def options_struct(endpoint)
    literals = []
    return literals unless endpoint.params?

    endpoint.all_params.dup.sort_by { |e| e.go_field_name }.each do |field|
      literal = field.description.empty? ? '' : "\n#{format_go_comment(field.description)}\n"
      literal += "#{field.go_protofield_name} #{field.ptr_go_type}"
      literal += "`bson:\"#{field.sql_identifier}\" json:\"#{field.identifier}\" sql:\"#{field.sql_identifier}\"`"
      literals << literal
    end
    comment = format_go_comment("#{endpoint.go_model_name}Options are options for API requests.")
    "\n#{comment}\ntype #{endpoint.go_model_name}Options struct {#{literals.sort.join("\n")}}"
  end
end
