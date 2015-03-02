/**
 * Created by Yinxiong on 2014-11-05.
 */

define(['kernel', 'angular', './module'],
function( core, ng, service){

    service.factory('Websocket', function( $websocket ){

        var Socket = core.Class.extend({
            init: function( events, reconnectOnLoss ){
                this._reconnectOnLoss = ng.isUndefined(reconnectOnLoss) ? true : reconnectOnLoss;
                this.events = ng.extend({
                    message: {},
                    open: ng.noop,
                    error: ng.noop,
                    close: ng.noop,
                    giveup: ng.noop
                }, events);
            },
            _isListen: false,
            _isConnected: false,
            _retryTimes: 0,
            listen: function(){
                var self = this;

                if( this._isListen ){
                    return;
                }
                this._isListen = true;

                this.socket = $websocket('ws://'+location.host+'/socket/');

                this.socket.onMessage(function( message ){
                    var data = JSON.parse(message.data);
                    if( data.method in self.events.message ){
                        self.events.message[data.method]( data.data, message )
                    }
                });

                this.socket.onOpen(function(e){
                    self.events.open.call(self, e);
                    self._retryTimes = 0;
                    self._isConnected = true;
                });

                this.socket.onClose(function( e ){
                    self.events.close.call(self, e);
                    self._reconnectOnLoss && self.reconnect();
                    self._isConnected = false;
                });

                this.socket.onError(function( e ){
                    if( !self._isConnected ) return;

                    self.events.error.call(self, e);
                    self._reconnectOnLoss && self.reconnect();
                });

                return this;
            },
            reconnect: function(){
                if( this._retryTimes >= 3 ){
                    this.events.giveup();
                    this.disconnect(true);
                    this._isListen = false;
                    return;
                }
                this._retryTimes++;
                console.warn('websocket disconnect, retry', this._retryTimes);
                this.socket.reconnect();
            },
            disconnect: function( force ){
                if( force ){
                    this._reconnectOnLoss = false;
                }
                this.socket.close( force );
            },
            emit: function( method, data ){
                return this.send({
                    method: method,
                    data: data
                })
            },
            broadcast: function( method, data ){
                return this.emit('broadcast', {
                    method: method,
                    data: data
                });
            },
            on: function( method, fn ){
                if( typeof method != 'string' ) return;

                if( (/^(error|close|open|giveup)$/).test(method) ){
                    this.events[method] = fn;
                }else{
                    this.events.message[method] = fn;
                }
            },
            send: function( data ){
                return this.socket.send(data)
            }
        });

        return function( events, reconnect ){
            return new Socket( events, reconnect )
        }
    });
});