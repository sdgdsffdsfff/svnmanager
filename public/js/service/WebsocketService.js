/**
 * Created by Yinxiong on 2014-11-05.
 */

define(['kernel', 'angular', './module'],
function( core, ng, service){

    service.factory('WebsocketService', function(){

        var Socket = core.Class.extend({
            _onMsg: [],
            _onCls: [],
            _onErr: [],
            init: function( url, reconnectOnLoss ){
                this.url = url || 'ws://'+location.host+'/socket/';
                this._reconnectOnLoss = ng.isUndefined(reconnectOnLoss) ? true : reconnectOnLoss;
            },
            on: function(type, fn){

                switch(type) {
                    case 'message':
                        this._onMsg.push(fn);
                        break;
                    case 'close':
                        this._onCls.push(fn);
                        break;
                    case 'error':
                        this._onErr.push(fn);
                        break;
                }

                return this;
            },
            listen: function(){
                var self = this;

                this.ws.addEventListener("message", function(e) {
                    self._onMsg.map(function(v){
                        v.call(self, e)
                    });
                });

                this.ws.addEventListener("close", function(e) {
                    self._onCls.map(function(v){
                        v.call(self, e)
                    });
                    self._reconnectOnLoss && self.connect();
                });

                this.ws.addEventListener("error", function(e) {
                    self._onErr.map(function(v){
                        v.call(self, e)
                    });
                });
            },
            connect: function(){
                this.ws = new WebSocket(this.url);
                this.listen();
            },
            disconnect: function( truly ){
                if( truly ){
                    this._reconnectOnLoss = false
                }
                this.ws.close();
            },
            send: function( data ){
                data = ng.isObject(data) ? JSON.stringify(data) : data;
                this.ws.send(data)
            }
        });

        return {
            socket: function( url ){
                return new Socket( url )
            }
        }
    });
});