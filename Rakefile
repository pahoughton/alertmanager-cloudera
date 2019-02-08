# 2019-02-07 (cc) <paul4hough@gmail.com>
# y = rake recurses down (. .. ../..:)

$runstart = Time.now

at_exit {
  runtime = Time.at(Time.now - $runstart).utc.strftime("%H:%M:%S.%3N")
  puts "run time: #{runtime}"
}

task :default do
  sh 'rake --tasks'
  exit 1
end

task :build_static do
  sh 'go build -mod=vendor ' + \
     '-tags netgo -ldflags \'-w -extldflags "-static"\''
end

task :build do
  sh 'go build -mod=vendor'
end

task :check => [:build] do
  sh './cloudera-amgr-alert --debug cloudera-alert.json'
end
