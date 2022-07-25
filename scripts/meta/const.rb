# frozen_string_literal: true

# Const is constant data to be used over any number of classes and modules used
# by the metaprogrammer.
module Const
  GEN_MSG = "\n// * This is a generated file, do not edit\n"
  MODEL_FILENAME = 'models.go'
  PARENT_DIR = File.expand_path('..', Dir.pwd)
  RETURN_ERR	= 'if err != nil { return err }'
  SERIAL_PKG = 'github.com/alpine-hodler/driver/internal/serial'
	TYPE_FILENAME = 'types.go'

  def self.go_fmt(filename)
    `/go/bin/goimports -w #{filename}`
  end

  def self.import(package)
    "\nimport \"#{package}\";"
  end

  def self.package(name)
    "package #{name}"
  end
end
