# 2019-02-07 (cc) <paul4hough@gmail.com>
# y = rake recurses down (. .. ../..:)

$runstart = Time.now

at_exit {
  runtime = Time.at(Time.now - $runstart).utc.strftime("%H:%M:%S.%3N")
  puts "run time: #{runtime}"
}

app = 'alertmanager-cloudera'
version = File.open('VERSION', &:readline).chomp

task :default do
  sh 'rake --tasks'
  exit 1
end

desc 'lint'
task :yamllint do
  sh "yamllint -f parsable .travis.yml .gitlab-ci.yml config"
end

desc 'yamllint && go test -v ./...'
task :test => [:yamllint] do
  sh 'go test -v ./...'
end
desc 'dev go build'
task :build do
  sh 'go build -mod=vendor'
end

desc 'test && static'
task :build_static => [:test] do
  require 'git'
  git = Git.open('.')

  branch = git.branch
  commit = git.gcommit('HEAD').sha
  version = File.open('VERSION', &:readline).chomp
  tag = git.tags[-1]

  sh 'go build -mod=vendor ' + \
     "-tags netgo -ldflags '" +\
     "-X main.Version=#{version} " +\
     "-X main.Branch=#{branch} " +\
     "-X main.Revision=#{commit} " +\
     "-X main.BuildDate=#{Time.now.strftime("%Y-%m-%d.%H:%M")} " +\
     "-w -extldflags -static'"
end

desc "create #{app}-#{version}.amd64.tar.gz"
task :release => [:build_static] do
  require 'git'
  git = Git.open('.')

  branch = git.branch
  commit = git.gcommit('HEAD').sha
  tag = git.tags[-1]
  if tag.sha != commit
    puts "rev not tagged"
    exit 1
  end
  if tag.name != "v#{version}"
    puts "tag '#{tag.name}' != 'v#{version}' VERSION file "
    exit 1
  end
  puts "app: #{app}"
  puts "branch: #{branch}"
  puts "commit: #{commit}"
  puts "version: #{version}"
  modified = false
  git.status.each do |f|
    if f.type || f.untracked
      mod = f.untracked ? "U" : f.type
      puts "#{mod} " + f.path
      modified = true
    end
  end
  if modified
    puts "modified or untracked files exists"
    exit 1
  end
  sh "mkdir #{app}-#{version}.amd64"
  sh "cp #{app} README.md VERSION COPYING #{app}-#{version}.amd64"
  sh "tar czf #{app}-#{version}.amd64.tar.gz #{app}-#{version}.amd64"
  sh "tar tzf #{app}-#{version}.amd64.tar.gz"
end

task :travis do
  sh 'go test -v ./...'
end
