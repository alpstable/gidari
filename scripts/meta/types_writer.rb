# frozen_string_literal: true

# TypesWriter is a module that contains the methods used to write the type.go files for web API packages.
module TypesWriter
  def self.apis(types)
    tree = (proc { Hash.new { |hash, key| hash[key] = [] } }).call
    types.each { |t| tree[t.api] << t }
    tree
  end

  def self.stringer(enum)
    logic = "\n// String will convert a #{enum.go_type_name} into a string.\n"
    logic += "func (#{enum.go_var_name} *#{enum.go_type_name}) String() string {;"
    logic += "if #{enum.go_var_name} != nil { return string(*#{enum.go_var_name}) };"
    logic += "return \"\"; };\n"
    unless enum.pluralize.nil?
      logic += "\n// String will convert a slice of #{enum.go_type_name} into a CSV.\n"
      logic += "func (#{enum.pluralize_var} *#{enum.pluralize}) String() string {;"
      logic += 'var str string;'
      logic += "if #{enum.pluralize_var} != nil {;"
      logic += 'slice := []string{};'
      logic += "for _, val := range *#{enum.pluralize_var} {;"
      logic += 'slice = append(slice, val.String());'
      logic += '};'
      logic += 'str = strings.Join(slice, ",");'
      logic += '};'
      logic += "return str;};\n"
    end
    logic
  end

  def self.enum_structure_contents(enum)
    enum.values.dup.map do |field|
      comment = field.go_comment.nil? ? '' : "\n#{field.go_comment}\n"
      "#{comment}#{enum.go_type_name + field.go_field_name} #{enum.go_type_name}" \
      " = \"#{field.identifier}\"\n"
    end.join('')
  end

  def self.enum_structure(api_type)
    api_type.enums.dup.map do |enum|
      structure = "\n#{enum.go_comment}\n"
      structure += "type #{enum.go_type_name} string;"
      structure += "type #{enum.pluralize} []#{enum.go_type_name};" unless enum.pluralize.nil?
      structure += "const (#{enum_structure_contents(enum)});"
      structure += stringer(enum)
      structure
    end
  end

  def self.enum_structures(api_types)
    api_types.dup.map do |api_type|
      enum_structure(api_type)
    end.flatten.sort.join("\n")
  end

  def self.write_file(file_, api, types)
    file_.write(Const.package(api))
    file_.write(Const::GEN_MSG)
    file_.write(enum_structures(types))
  end

  def self.write(types)
    apis(types).each do |api, api_types|
      path = Pathname.new(PARENT_DIR).join(api)
      Dir.chdir(path.to_s) do
        File.open(Const::TYPE_FILENAME, 'w') do |f|
          write_file(f, api, api_types)
        end
        Const.go_fmt(Const::TYPE_FILENAME)
      end
    end
  end
end
