/*
open server
  on request > create new organism

organism is an object with everything it needs to survive in the
ecosystem, including but not limited to:
  responding to its request
  being read by functions
  running functions

*/


"use strict";

const events = require('events');
const fs = require('fs');
const http = require('http');


//var server    = http.createServer(router);
var server    = https.createServer(conception);


function conception(req, res) {
	console.log('fertilizing egg')
	var organism  = new Organism(req, res);
	console.log('new Organism: \n' + organism);
	//get body
	req.on('data', function(chunk) {
		organism.body.push(chunk);
		if(organism.body.length > 1e7) {
          res.writeHead(413, 'Request Entity Too Large');
          res.end('file too large');
					organism.kill();
        }
	}).on('end', function() {
		organism.body = Buffer.concat(organism.body).toString();
    console.log(organism.body);
    if(organism.body){
      organism.body = JSON.parse(organism.body);
    }
    //conception of organism complete, being journey
    // api.ear7h.net only handles single level paths to api calls
    //ie. api.ear7h.net/login, api.ear7h.net/isUniqueUsername
    console.log(req.headers.host);
		function wow (){
			organism.emit('death');
		}
    organism.call(req.url, wow);
	});
}

function parseCookie(cookString) {
	if(cookString){} else {
		return null;
	}
	var cookArray = cookString.split("; ");
	var allCookies = {};
	for (i = 0; i < cookArray.length; i++){
		var thisCookie = cookArray[i].split("=");
		allCookies[thisCookie[0]] = JSON.parse(thisCookie[1]);
	}
	console.log("cookie" + JSON.stringify(allCookies));
	return allCookies;
}

class Organism extends events {
	constructor(req, res) {
		super();
		this.req 					= req;
		this.res 					= res;
		this.ip						= req.headers['x-forwarded-for'] ||
			req.connection.remoteAddress ||
			req.socket.remoteAddress ||
			req.connection.socket.remoteAddress;
		this.result		= {
			section		: null,
			state 		: null,
			mime			: null,
			data			: null
		};
		this.body 				= [];
		this.httpHeader 	= {
			code			: 200,
			fields		: {
				'content-type': 'application/json',
			},
		};
		//console.log(this);
		//start declarations of events
		this.on('death', this.death);
		//end declations
	}
	//events			:new events(); //set by constructor
	death       	(){
		console.log('organism died');
		this.resSend();
	}
	kill        	(){
		console.log('killing organism');
		this.emit('death');
	}
	resSend     	(){
		console.log('resSend');
		var header = this.httpHeader;
		console.log('RESULT:\n' + JSON.stringify(this.result) + '\n \n');
		this.res.writeHead(header.code, header.fields);
		this.res.write(JSON.stringify(this.result));
		this.res.end()
		console.log('\n API END \n');
	}
	setCookie			(string){
		this.httpHeader.fields['set-cookie'] = string;
		console.log('httpHeader.set-cookie = ' + this.httpHeader.fields['set-cookie']);
	}
	getCookie		(){
		 parseCookie(req.headers.cookie);
	 }
	 call        (jsFile, callback){
		try{
			console.log('calling: ' + jsFile + '.js');
			require('.' + jsFile + '.js').main(this, callback);

		} catch (e){
			console.log('organism call error: ' + e.stack);
			this.result.data = 'invalid';
			this.kill();
		}
	}

}
//api server will send things in json format, exclusively
server.listen(80);
