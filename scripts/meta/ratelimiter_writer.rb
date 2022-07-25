# frozen_string_literal: true

# TypRatelimitWriteresWriter is a module that contains the methods used to write the ratelimiter.go files for web API
# packages.
module RatelimitWriter
  def self.apis(ratelimiters)
    tree = (proc { Hash.new { |hash, key| hash[key] = [] } }).call
    ratelimiters.each do |rl|
			next if rl.endpoint.nil?

      tree[rl.endpoint.api] << rl
    end
    tree
  end

  def self.consts(ratelimiters)
    consts = ['_ ratelimiter = iota'] | ratelimiters.dup.map { |rl| rl.endpoint.go_rl_const }.sort
    "const(#{consts.join(';')});"
  end

  def self.runlimiter_type
    "\ntype ratelimiter uint8;\n"
  end

  def self.runlimiter_var(ratelimiters)
    "var ratelimiters = [uint8(#{ratelimiters.length + 1})]*rate.Limiter{};"
  end

  def self.init ratelimiters
    consts = ratelimiters.dup.map { |rl| "ratelimiters[#{rl.endpoint.go_rl_const}] = nil;" }.sort
    "func init() {#{consts.join(';')}};\n\n"
  end

  def self.fn
		logic = "// getRateLimiter will load the rate limiter for a specific request, lazy loaded.\n"
    logic += 'func getRateLimiter(rl ratelimiter, b int) *rate.Limiter {;'
    logic += 'if ratelimiters[rl] == nil {;'
    logic += 'ratelimiters[rl] = rate.NewLimiter(rate.Every(1*time.Second),b);'
    logic += '};'
    logic += 'return ratelimiters[rl];'
    logic += "};\n\n"
    logic
  end

  def self.write(ratelimiters)
    apis(ratelimiters).each do |api, ratelimiters|
      path = Pathname.new(PARENT_DIR).join(api)
      Dir.chdir(path.to_s) do
        File.open('ratelimiter.go', 'w') do |f|
          f.write("package #{api}")
          f.write("\nimport \"golang.org/x/time/rate\";")
          f.write("\nimport \"time\";")
          f.write(GEN_MSG)
          f.write(runlimiter_type)
          f.write(consts(ratelimiters))
          f.write(runlimiter_var(ratelimiters))
          f.write(init(ratelimiters))
          f.write(fn)
        end
        `/go/bin/goimports -w ratelimiter.go`
      end
    end
  end
end
