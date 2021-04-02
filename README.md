# Welcome to Revel

A high-productivity web framework for the [Go language](http://www.golang.org/).


### Start the web server:

   revel run myapp

### Go to http://localhost:9000/ and you'll see:

    "It works"

## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites


## Help

* The [Getting Started with Revel](http://revel.github.io/tutorial/gettingstarted.html).
* The [Revel guides](http://revel.github.io/manual/index.html).
* The [Revel sample apps](http://revel.github.io/examples/index.html).
* The [API documentation](https://godoc.org/github.com/revel/revel).



## Makefile

# Copyright 2017 Google Inc. All rights reserved.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to writing, software distributed
# under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
# CONDITIONS OF ANY KIND, either express or implied.
#
# See the License for the specific language governing permissions and
# limitations under the License.

build:
	#GOOS=linux CGO_ENABLED=0 go build -o app
	#docker build -t campoy/say .
	#1. first remove unzips folder and its content
	rm -rf unzips
	rm -rf certs-162.252.80.136
	rm -f eventvisorHQ.tar.gz
	#2. create folder
	mkdir unzips
	mkdir -p uploads/conference/poster
	mkdir -p uploads/conference/thumbnail
	mkdir -p uploads/session/poster
	mkdir -p uploads/session/thumbnail
	mkdir -p uploads/user/profile
	mkdir certs-162.252.80.136
	#3. package revel command
	#revel package -m prod
	#CGO_ENABLED=0 revel package github.com/najamsk/eventvisor/eventvisorHQ -s
	CGO_ENABLED=0 revel package github.com/najamsk/eventvisor/eventvisorHQ -m prod
	cp  /home/najam/go/src/github.com/najamsk/DeploymentsSecure/eventvisor/eventvisorHQ/.env .
	cp -r /home/najam/go/src/github.com/najamsk/DeploymentsSecure/eventvisor/eventvisorHQ/certs-162.252.80.136/* certs-162.252.80.136
	chmod 0600 certs-162.252.80.136/client.root.key
	. ./.env
	#4. copy package tar file to unzips folder
	#cp eventvisorHQ.tar.gz unzips
	#cp Dockerfile unzips
	#5. untar to unzips folder 
	#tar -zxvf eventvisorHQ.tar.gz -C unzips
	#6. remove tar file 
	#rm -f unzips/eventvisorHQ.tar.gz
	#rm -f eventvisorHQ.tar.gz
	
	#7. run revel
	#unzips/run.sh

image:
	#docker build -t eventvisor/hq unzips/.
	docker build -t eventvisorhq .

export:
	docker save eventvisorhq > hqdocker.tar

clean: 
	rm -rf unzips
	rm .env
	rm -rf certs-162.252.80.136
	rm -f eventvisorHQ.tar.gz
	rm -f hqdocker.tar

container:
	#docker run --rm -it --cap-add=SYS_PTRACE --security-opt seccomp=unconfined --name eventvisorHQ -p 9000:9000 eventvisorhq 
	#docker run --rm -it --name eventvisorHQ -v uploads:/app/src/github.com/najamsk/eventvisor/eventvisorHQ/uploads -p 9000:9000 eventvisorhq 
	docker run --rm -it --name eventvisorHQ -v "$(CURDIR)"/uploads:/app/src/github.com/najamsk/eventvisor/eventvisorHQ/uploads --env-file ./.env -p 9000:9000 eventvisorhq

run:
	#GOOS=linux CGO_ENABLED=0 go build -o app
	#docker build -t campoy/say .
	#1. first remove unzips folder and its content
	rm -rf unzips
	rm -f eventvisorHQ.tar.gz
	#2. create folder
	mkdir unzips
	#3. package revel command
	revel package -m prod
	#4. copy package tar file to unzips folder
	cp eventvisorHQ.tar.gz unzips
	#rm -f eventvisorHQ.tar.gz
	#5. untar to unzips folder 
	tar -zxvf eventvisorHQ.tar.gz -C unzips
	#6. remove tar file 
	rm -f eventvisorHQ.tar.gz
	#7. run revel
	unzips/run.sh
	
dev:
	#GOOS=linux CGO_ENABLED=0 go build -o app
	#docker build -t campoy/say .
	#1. first remove unzips folder and its content
	rm -rf unzips
	#rm -rf uploads
	rm -f eventvisorHQ.tar.gz
	#2. create folder
	mkdir unzips
	# -p flag slient the already exist error and excute rest of the script
	mkdir -p uploads
	#3. package revel command
	revel package
	#4. copy package tar file to unzips folder
	cp eventvisorHQ.tar.gz unzips
	#rm -f eventvisorHQ.tar.gz
	#5. untar to unzips folder 
	tar -zxvf eventvisorHQ.tar.gz -C unzips
	#6. remove tar file 
	rm -f eventvisorHQ.tar.gz
	#7. run revel
	unzips/run.sh



    ## Minification

    ### CSS
    minify -o secure-all.css site.css font-awesome.min.css responsive.css dropzone.css eventvisor.css bootstrap-datetimepicker.min.css bootstrap-datepicker.css daterangepicker.css datepicker.css timepicker.css datetime.css datepick.css bootstrap-glyphicons.css

	minify -o public-all.css site.css responsive.css custom.css

    ### JS
    minify -o secure_all.js jquery-3.2.1.min.js modernizr.custom.js jquery-ui.min.js popper.min.js bootstrap.min.js moment.min.js bootstrap-datetimepicker.min.js bootstrap-datepicker.js daterangepicker.min.js bootstrap-timepicker.min.js moment-timezone-with-data.min.js dropzone.min.js scripts.js

	minify -o public_all.js jquery-3.2.1.min.js modernizr.custom.js jquery-ui.min.js popper.min.js bootstrap.min.js moment.min.js bootstrap-datetimepicker.min.js bootstrap-datepicker.js daterangepicker.min.js bootstrap-timepicker.min.js moment-timezone-with-data.min.js scripts.js


	##sayyam Minification

    ### CSS
    minify -o secure-all.css font-awesome.min.css dropzone.css eventvisor.css bootstrap-datetimepicker.min.css bootstrap-datepicker.css daterangepicker.css datepicker.css timepicker.css datetime.css datepick.css bootstrap-glyphicons.css

	minify -o public-all.css custom.css

	minify -o common-all.css site.css responsive.css


    ### JS
    minify -o secure_all.js dropzone.min.js 

	minify -o common_all.js jquery-3.2.1.min.js modernizr.custom.js jquery-ui.min.js popper.min.js bootstrap.min.js moment.min.js bootstrap-datetimepicker.min.js bootstrap-datepicker.js daterangepicker.min.js bootstrap-timepicker.min.js moment-timezone-with-data.min.js scripts.js