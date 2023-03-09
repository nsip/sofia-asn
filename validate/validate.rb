## Validation of MRAC JSON-LD output. 
## presupposes that it is run in a folder containing the JSON-LD output for MRAC, as *.json files,
## and the ScOT JSON-LD file downloaded from https://scot.curriculum.edu.au, as scot.jsonld
## Validation checks: whether all JSON keys are namespaced as JSON-LD keys; whether the
## gem:hasChild/gem:isChildOf pairs match; whether cross-references to other parts of MRAC
## are pointing to the right kind of referent; and whether ScOT cross-references are legal.

def scot_id_parse(json, ids)
  case json
  when Hash
    if (k = id(json)) &&
       %r{^http://vocabulary.curriculum.edu.au/scot/}.match?(k)
      ids[k] = true
    end
    json.each_value do |m|
      scot_id_parse(m, ids)
    end
  when Array
    json.each { |j| scot_id_parse(j, ids) }
  end
  ids
end

def asn_id_parse(json, ids)
  if (k = id(json)) && json['asn:statementLabel'] # no label, it's a root of tree, not a real ID
    ids[k] and warn "#{@f}: id_parse: #{k} is duplicated beween #{@ids[k]} and #{@f}!"
    ids[k] = { file: @f, label: json['asn:statementLabel'] }
  end
  json['gem:hasChild']&.each do |m|
    asn_id_parse(m, ids)
  end
  ids
end

def child_parent(json)
  json['gem:hasChild']&.each do |j|
    j['gem:isChildOf']&.each do |k|
      id(json) or next
      (id(json) == id(k)) and next
      warn "#{@f}: child_parent: #{id(j)} mismatches parent and child!"
    end
    child_parent(j)
  end
end

def id(x)
  x['@id'] || x['@Id']
end

def namespaced_keys(json, key)
  case json
  when Hash
    json.each_key do |k|
      namespaced_keys(json[k], k)
      key == '@context' and next
      !/\s/.match?(k) && /\S:\S|^@/.match?(k) and next
      warn "#{@f}: namespaced_keys: #{k} is an illegal key!"
    end
  when Array
    json.each { |j| namespaced_keys(j, nil) }
  end
end

def exist_link(k, lbl)
  @ids[k] or
    warn "#{@f}: link_types: #{k} destination of #{lbl} is not a defined key!"
  @ids[k]
end

def link_types(json)
  case json
  when Hash
    # require "debug"; binding.b if (@f == "la-English.json.json")
    if (k = json.dig('asn:skillEmbodied', '@id')) && exist_link(k, 'asn:skillEmbodied')
      /^gc/.match?(@ids[k][:file]) or
        warn "#{@f}: link_types: #{k} destination of asn:skillEmbodied is not a general capability!"
    end
    if (k = json.dig('asn:crossSubjectReference', '@id')) && exist_link(k, 'asn:crossSubjectReference')
      /^ccp/.match?(@ids[k][:file]) or
        warn "#{@f}: link_types: #{k} destination of asn:crossSubjectReference is not a cross-curriculum priority!"
    end
    if (k = json.dig('asn:hasLevel', '@id')) && exist_link(k, 'asn:hasLevel')
      (@ids[k][:label] == 'Level') or
        warn "#{@f}: link_types: #{k} destination of asn:hasLevel is not an achievement level!"
    end
    if (k = json.dig('dc:relation', '@id')) && exist_link(k, 'dc:relation')
      /^la/.match?(@ids[k][:file]) or
        warn "#{@f}: link_types: #{k} destination of dc:relation is not a curriculum statement!"
    end
    json.each_value do |m|
      link_types(m)
    end
  when Array
    json.each { |j| link_types(j) }
  end
end

def scot(json)
  case json
  when Hash
    (s = json['asn:conceptTerm']) and scot1(s)
    json.each_value do |m|
      scot(m)
    end
  when Array
    json.each { |j| scot(j) }
  end
end

def scot1(terms)
  ret = if terms.is_a?(Array) then terms.map { |x| x['@id'] }
        else
          [terms['@id']]
        end
  ret.each do |k|
    k.nil? and warn "#{@f}: scot: #{terms} destination of asn:conceptTerm is undefined!"
    @scotids[k] or
      warn "#{@f}: scot: #{k} destination of asn:conceptTerm is undefined!"
  end
end

def validate_file(json)
  child_parent(json)
  namespaced_keys(json, nil)
  link_types(json)
  scot(json)
end

def read_file(f)
  JSON.parse(File.read(f))
end

def read_scot
  scot = JSON.parse(File.read('scot.jsonld'))
  @scotids = scot_id_parse(scot, {})
end

def validate
  read_scot
  @ids = {}
  Dir.glob('*.json') do |f|
    @f = f
    a = read_file(f)
    @ids = asn_id_parse(a, @ids)
  end
  Dir.glob('*.json') do |f|
    @f = f
    a = read_file(f)
    validate_file(a)
  end
end

validate
