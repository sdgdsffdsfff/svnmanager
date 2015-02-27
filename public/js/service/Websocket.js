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
                    error: ng.noop,
                    close: ng.noop
                }, events);
            },
            _isListen: false,
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

                this.socket.onClose(function( e ){
                    self.events.close.call(self, e);
                    self._reconnectOnLoss && self.socket.reconnect();
                });

                this.socket.onError(function( e ){
                    self.events.error.call(self, e);
                });

                return this;
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
            on: function( method, fn ){
                if( typeof method != 'string' ) return;

                if( method == 'error' ){
                    this.events['error'] = fn;
                } else if( method == 'close' ){
                    this.events['close'] = fn;
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