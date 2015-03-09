/**
 * Created by Yinxiong on 2014-11-05.
 */

define([
'kernel',
'angular',
'./module',
'ui/Toast',
'service/Websocket'
],
function( core, ng, service, Toast){

    service.factory('SocketInstance', function( Websocket ) {
        var socket = Websocket();
        socket.scope = null;
        socket.setScope = function( s ){
            this.scope = s;
        };
        socket.on('notify', function( msg ){
            Toast.makeText( msg ).show();
        });
        return socket.listen();
    })
});