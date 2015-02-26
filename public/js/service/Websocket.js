/**
 * Created by Yinxiong on 2014-11-05.
 */

define(['kernel', 'angular', './module'],
function( core, ng, service){

    service.factory('Websocket', function( $websocket ){

        var Socket = core.Class.extend({
            init: function( events, reconnectOnLoss ){
                this.socket = $websocket('ws://'+location.host+'/socket/');
                this._reconnectOnLoss = ng.isUndefined(reconnectOnLoss) ? true : reconnectOnLoss;
                this.events = ng.extend({
                    message: {},
                    error: {},
                    close: {}
                }, events);
                this.listen();
            },
            listen: function(){
                var self = this;

                this.socket.onMessage(function( message ){
                    var data = JSON.parse(message.data);
                    if( data.method in self.events.message ){
                        self.events.message[data.method]( data.data, message )
                    }
                });

                this.socket.onClose(function( e ){
                    for( var i in self.events.close ){
                        self.events.close[i](e)
                    }
                    self._reconnectOnLoss && self.socket.reconnect();
                });

                this.socket.onError(function( e ){
                    for( var i in self.events.error ){
                        self.events.error[i](e)
                    }
                });
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
            send: function( data ){
                return this.socket.send(data)
            }
        });

        return function( events, reconnect ){
            return new Socket( events, reconnect )
        }
    });
});