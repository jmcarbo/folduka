#!/usr/bin/env ruby

require "openssl"
require "time"

begin
  require "origami"
rescue LoadError
  abort "origami not installed: gem install origami"
end
include Origami

CERT_FILE = "certificate.crt"
KEY_FILE = "private-key.pem"

input_files = ARGV

if input_files.empty?
  abort "Usage: sign-pdf input.pdf [...]"
end

key = OpenSSL::PKey::RSA.new(File.read(KEY_FILE))
cert = OpenSSL::X509::Certificate.new(File.read(CERT_FILE))

input_files.each do |file|
  output_filename = file.dup.insert(file.rindex("."), "_signed")

  pdf = PDF.read(file)
  page = pdf.get_page(1)

  width = 200.0
  height = 50.0
  x = page.MediaBox[2].to_f - width - height
  y = height
  size = 8

  now = Time.now

  text_annotation = Annotation::AppearanceStream.new
  text_annotation.Type = Origami::Name.new("XObject")
  text_annotation.Resources = Resources.new
  text_annotation.Resources.ProcSet = [Origami::Name.new("Text")]
  text_annotation.set_indirect(true)
  text_annotation.Matrix = [ 1, 0, 0, 1, 0, 0 ]
  text_annotation.BBox = [ 0, 0, width, height ]
  text_annotation.write("Signed at #{now.iso8601}", x: size, y: (height/2)-(size/2), size: size)

  # Add signature annotation (so it becomes visibles in PDF document)
  signature_annotation = Annotation::Widget::Signature.new
  signature_annotation.Rect = Rectangle[llx: x, lly: y+height, urx: x+width, ury: y]
  signature_annotation.F = Annotation::Flags::PRINT
  signature_annotation.set_normal_appearance(text_annotation)

  page.add_annotation(signature_annotation)

  # Sign the PDF with the specified keys
  pdf.sign(cert, key,
    method: "adbe.pkcs7.sha1",
    annotation: signature_annotation,
    location: "Helsinki",
    contact: "contact@kiskolabs.com",
    reason: "Proof of Concept"
  )

  # Save the resulting file
  pdf.save(output_filename)
end
