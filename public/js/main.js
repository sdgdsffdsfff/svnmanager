var components = '../components';

requirejs.config({
    paths: {
        kernel: 'core/kernel',
        jquery: "utils/jquery",
        react: components + '/react/react-with-addons',
        TweenMax: components + '/greensock/src/minified/TweenMax.min',
        angular: components + '/angularjs/angular.min',
        ngSanitize: components + '/angular-sanitize/angular-sanitize.min',
        ngWebSocket: components + '/angular-websocket/angular-websocket.min',
        bootstrap: components + '/bootstrap/dist/js/bootstrap.min',
        socket: components + '/socket.io-client/socket.io',
        moment: components + '/moment/min/moment.min',
        text: components + '/requirejs-text/text'
    },
    map :{
        '*' : {
            'css': components + '/require-css/css.min'
        }
    },
    shim: {
        angular: {
            exports: 'angular'
        },
        ngSanitize: {
            deps: ['angular']
        },
        ngWebSocket: {
            deps: ['angular']
        },
        moment: {
            exports: 'moment'
        }
    }
});

require(['init', 'ui/bootstrap']);