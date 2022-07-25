# frozen_string_literal: true

# Comment is a module for methods to maninipulate comments for a specific
# programming language
module Comment
  def self.split_comment_by_char(delimiter, comment, count)
    line =	delimiter
    split = []
    comment.split(' ').each do |word|
      if (line + word).length >= count
        split << (line += "\n").rstrip
        line = delimiter
      end
      line = line + ' ' + word
    end
    split << line.rstrip
  end

  def self.split_go_comment_by_char(comment, count = 120)
    Comment.split_comment_by_char('//', comment, count)
  end

  def self.split_graphql_comment_by_char(comment, count = 120)
    ['"""', Comment.split_comment_by_char('', comment, count), '"""']
  end

  def self.u_format_go_comment(comment, count = 120)
    split_go_comment_by_char(comment, count).join("\n")
  end

  def format_go_comment(comment, count = 120)
    Comment.split_go_comment_by_char(comment, count).join("\n")
  end

  def format_go_comment_with_source(comment, source, count = 120)
    fmt = format_go_comment(comment, count)
    fmt += "\n//\n"
    fmt += "// source: #{source}"

    fmt
  end

  def format_graphql_comment(comment, count = 120)
    Comment.split_graphql_comment_by_char(comment, count).join("\n")
  end
end
