/**
 * Created by Yinxiong on 2014-10-10.
 */
define(function(){
    var core = {};
    var root = this;
    var $window = $(window);

    core.win = $window;
    core.body = null;
    $(function(){
        core.body = $('body');
    });

    core.Class = core.Class = function(){
        var initializing = false, fnTest = /xyz/.test(function(){xyz;}) ? /\b_super\b/ : /.*/;

        // The base Class implementation (does nothing)
        var Class = function(){};

        // Create a new Class that inherits from this class
        Class.extend = function(prop, statics) {
            var _super = this.prototype;

            // Instantiate a base class (but only create the instance,
            // don't run the init constructor)
            initializing = true;
            var prototype = new this();
            initializing = false;

            // Copy the properties over onto the new prototype
            for (var name in prop) {
                // Check if we're overwriting an existing function
                prototype[name] = typeof prop[name] == "function" &&
                    typeof _super[name] == "function" && fnTest.test(prop[name]) ?
                    (function(name, fn){
                        return function() {
                            var tmp = this._super;

                            // Add a new ._super() method that is the same method
                            // but on the super-class
                            this._super = _super[name];

                            // The method only need to be bound temporarily, so we
                            // remove it when we're done executing
                            var ret = fn.apply(this, arguments);
                            this._super = tmp;

                            return ret;
                        };
                    })(name, prop[name]) :
                    prop[name];
            }

            // The dummy class constructor
            function Class() {
                // All construction is actually done in the init method
                if ( !initializing && this.init )
                    this.init.apply(this, arguments);
            }

            // Populate our constructed prototype object
            Class.prototype = prototype;

            // Enforce the constructor to be what we expect
            Class.prototype.constructor = Class;

            // And make this class extendable
            Class.extend = arguments.callee;

            for(var i in statics){
                Class[i] = statics[i];
            }

            if( Class.init ){
                Class.init.call(Class);
            }

            return Class;
        };

        return Class;
    }();

    core.Service = function(member){
        if( typeof member !== 'function' ) return null;

        return core.Class.extend(member());
    };

    core.extend = $.extend;

    //register namespace
    core.error = [];
    core.register = function(e, c) {
        var g = e.split(".");
        var f = core;
        var b = core.error;
        var d = null;
        while (d = g.shift()) {
            if (g.length) {
                if (f[d] === undefined) {
                    f[d] = {}
                }
                f = f[d]
            } else {
                if (f[d] === undefined) {
                    if( typeof c === 'string' && c === 'check' ){
                        core.console.log('Property undefined:'+e);
                        return $.noop;
                    }
                    try {
                        f[d] = c(core);
                    } catch(h) {
                        b.push(h)
                    }
                }else{
                    if( typeof c == 'string' && c == 'check' ){
                        return f[d]
                    }
                    console.log('redefined:'+e)
                }
            }
        }
    };

    core.getMethod = function(e){
        return core.register( e, 'check')
    };

    //延迟函数
    core.delay = function(fn, delay){
        return setTimeout(fn, delay || 0);
    };

    //图像加载
    core.loadImage = function(url, callback){
        var img = new Image();
        img.src = url;
        if (img.complete) {
            callback.call(img);
        }
        else {
            img.onload = function(){
                callback.call(img);
            };
            img.src = img.src
        }
    };

    //随机数
    core.random = function(target, start){
        if( typeof target == 'number' ){
            var count = target || 5;
            start = start || 1;
            return Math.round( Math.random() * Math.pow( 10,count ) ) * start;
        }
        if( (typeof target == 'object' || typeof target == 'string') && target.hasOwnProperty('length') ){
            var isString = typeof target == 'string';
            if( isString ){
                target = target.split('');
            }
            var obj = [], ex = {}, max = target.length || start || 10;
            var get = function(){
                var key = Math.round( Math.random()*max );
                return key in ex || key == max ? get() : ( ex[key] = true ) && key;
            };
            if(target.length == 0){
                for(var i = 0; i < max; i++){
                    obj[i] = get();
                }
            }else{
                for(var i = 0; i < max; i++){
                    obj.push(target[get()])
                }
            }
            if( isString ){
                return obj.join('');
            }
            return obj;
        }
    };

    /**
     * 修复当滚动和重置窗口大小时多次调用resize, scroll事件的bug
     * @param {[string, function], [object, function]} 事件名，配置项
     */
    core.xresize = function(name, options){
        var timer = null, isEnd = false;

        var setting = {
            //一绑定便执行
            init: true,
            //滚动停止后调用
            after: $.noop,
            //开始滚动时调用
            before: $.noop,
            delay: 100
        }

        if(arguments.length == 1){
            setting.after = name;
            name = 'scroll resize';
        }else{
            name = name || 'scroll resize';
            if( typeof options == 'function'){
                setting.after = options;
            }else{
                setting = $.extend(setting, options);
            }
        }

        core.win.bind(name, function(e){
            if(timer){
                clearTimeout(timer);
                timer = null;
            }
            if(isEnd == false){
                isEnd = true;
                setting.before.call(core.win, e)
            }
            timer = core.delay(function(){
                isEnd = false
                setting.after.call(core.win, e);
            }, setting.delay)
        });

        if(setting.init){
            core.win.trigger(name.split(' ')[0]);
        }

        return core.win;
    };

    core.xresize.unbind = function(name){
        core.win.unbind(name);
    };

    core.whenSeenElement = function(name, ele, callback){
        var name = 'resize.'+name+' scroll.'+name;
        core.xresize(name, function( e ){
            if( core.win.scrollTop()+core.win.height() < ele.offset().top ) return;
            core.xresize.unbind(name);
            (callback || $.noop)();
        })
    };

    /**
     * 固定布局
     * @param {jQuery object, object}  ele, options
     */
    core.positionFixed = function(){
        var id = 0, abs = Math.abs, ceil = Math.ceil;

        return function(ele, options){
            id++;
            var name = 'scroll.positionFixed' + id +' resize.positionFixed' + id;

            ele = ele.jquery ? ele : $(ele);

            var defaultValue = 'auto';

            options = $.extend({
                offset: {
                    left: defaultValue,
                    right: defaultValue,
                    top: defaultValue,
                    bottom: defaultValue
                },
                centerHorizontal: false,
                centerVertical: false,
                onPosition: $.noop,
                after: $.noop,
                before: $.noop
            }, options);

            ele.css(options.offset);
            if( options.centerHorizontal ){
                ele.css({
                    'left': '50%',
                    'marginLeft': -ele.outerWidth()/2
                });
            }
            if( options.centerVertial ){
                ele.css({
                    'top': '50%',
                    'marginTop': -ele.outerHeight()/2
                });
            }
            options.onPosition(ele);

            return name;
        }
    }();

    /**
     * 除了内容与触发手柄，点击其他位置都隐藏
     * @param {jQuery Object, jQuery Object}
     */
    core.clickAnyWhereHideButMe = function(){
        var pending = 0;

        return function (element, handle, callback) {
            if (arguments.length == 2) {
                callback = handle;
                handle = null;
            }

            var ele = element && element.jquery ? element[0] : element,
                h = handle && handle.jquery ? handle[0] : handle,
                name = 'mousedown.cawhb'+pending++;

            callback = callback || $.noop;

            $(document).bind(name, function (e) {
                if (( h ? e.target != h : true) && e.target != ele && !$.contains(ele, e.target)) {
                    callback(e, element);
                }
            });

            return name;
        }
    }();

    core.unbindDocumentEvent = function(name){
        $(document).unbind(name);
    };

    core.lazyLoad = function(options){
        options = $.extend({
            context : null,
            height : 0
        }, options);

        var win = core.win;
        var context = $(options.context);
        if(!context.length) return;
        var pageTop = function() {
            return document.documentElement.clientHeight
                + Math.max(document.documentElement.scrollTop, document.body.scrollTop)
                - options.height;
        };
        var imgLoad = function() {
            context.find('img[orgSrc]').each(function() {
                if($(this).offset().top <= pageTop() && $(this).is(':visible') ){
                    var orgSrc = this.getAttribute('orgSrc');
                    this.setAttribute('src', orgSrc);
                    this.removeAttribute('orgSrc');
                }
            });
        };
        win.bind('lazyload', imgLoad);
        core.xresize('scroll.lazyload',{
            after: imgLoad
        });
    };

    /**
     * like python range
     * TODO bug here
     */
    core.range = function(/*number*/ start, /*number*/ end, /*number -1 or > 0*/ step ){
        step = step == (-1 || step > 0) ? step : 1;
        if(arguments.length == 1){
            end = start;
            start = 0;
        }
        var arr = [], isReverse = step === -1;
        for(var i = start, step = isReverse ? 1 : step; i < end; i+=step){
            arr.push(i)
        }
        return isReverse ? arr.reverse() : arr;
    };

    /*!
     * jQuery Cookie Plugin v1.3.1
     * https://github.com/carhartl/jquery-cookie
     *
     * Copyright 2013 Klaus Hartl
     * Released under the MIT license
     */
    core.cookie = function(){
        var pluses = /\+/g;

        function raw(s) {
            return s;
        }

        function decoded(s) {
            return decodeURIComponent(s.replace(pluses, ' '));
        }

        function converted(s) {
            if (s.indexOf('"') === 0) {
                // This is a quoted cookie as according to RFC2068, unescape
                s = s.slice(1, -1).replace(/\\"/g, '"').replace(/\\\\/g, '\\');
            }
            try {
                return config.json ? JSON.parse(s) : s;
            } catch(er) {}
        }

        var config = cookie = function (key, value, options) {

            // write
            if (value !== undefined) {
                options = $.extend({}, config.defaults, options);

                if (typeof options.expires === 'number') {
                    var days = options.expires, t = options.expires = new Date();
                    t.setDate(t.getDate() + days);
                }

                value = config.json ? JSON.stringify(value) : String(value);

                return (document.cookie = [
                    encodeURIComponent(key), '=', config.raw ? value : encodeURIComponent(value),
                    options.expires ? '; expires=' + options.expires.toUTCString() : '', // use expires attribute, max-age is not supported by IE
                    options.path    ? '; path=' + options.path : '',
                    options.domain  ? '; domain=' + options.domain : '',
                    options.secure  ? '; secure' : ''
                ].join(''));
            }

            // read
            var decode = config.raw ? raw : decoded;
            var cookies = document.cookie.split('; ');
            var result = key ? undefined : {};
            for (var i = 0, l = cookies.length; i < l; i++) {
                var parts = cookies[i].split('=');
                var name = decode(parts.shift());
                var cookie = decode(parts.join('='));

                if (key && key === name) {
                    result = converted(cookie);
                    break;
                }

                if (!key) {
                    result[name] = converted(cookie);
                }
            }

            return result;
        };

        config.defaults = {};

        var removeCookie = function (key, options) {
            if (cookie(key) !== undefined) {
                cookie(key, '', $.extend(options, { expires: -1 }));
                return true;
            }
            return false;
        }

        cookie.removeCookie = removeCookie;

        return cookie
    }();

    core.go = function( url, isNewWindow ){
        var local = location, href = '';
        if( url == 'me'){
            href = local.href;
        }else if( /^#/.test(url) ){
            href = local.origin + local.pathname + url;
        }else{
            href = url;
        }
        if( !isNewWindow ){
            local.href = href;
        }else{
            window.open(href)
        }
        return local;
    };

    //伪ajax回调测试
    core.ajax = function(options){
        if( options.debug ){
            var deferred = new $.Deferred();
            var code = options.statusCode || 200;
            var delay = options.delay || 300;
            switch(code){
                case 200:
                    core.delay(function(){
                        deferred.resolve({
                            code: 1
                        });
                    }, delay)
                    break;
                case 403:
                case 404:
                case 500:
                    core.delay(function(){
                        deferred.reject({
                            code: 0
                        });
                    }, delay)
                    break;
            }
            return deferred;
        }else{
            return $.ajax(options);
        }
    };

    /**
     * HTML5 upload
     * @param url
     * @param data
     * @param options
     * @returns {XMLHttpRequest} 2.0
     */
    core.upload = function(url, data, options){
        options = $.extend({
            progress    : $.noop,
            load        : $.noop,
            error       : $.noop,
            abort       : $.noop
        }, options);

        var fd = new FormData();
        $.each(data, function( key, value ){
            fd.append(key, value);
        });

        var xhr = new XMLHttpRequest();

        xhr.addEventListener('load', function( e ){
            options.load.call(e, e.currentTarget.responseText);
        }, false);

        xhr.addEventListener('error', options.error, false);

        xhr.upload.addEventListener('progress', function(e){
            options.progress.call( e, Math.ceil((e.loaded / e.total) * 100), xhr );
        }, false);

        xhr.addEventListener('abort', options.abort, false);

        xhr.open('post', url, true);

        xhr.send(fd);

        return xhr;
    };

    core.fileSize = function( size ){
        return {
            M: Math.round(size/1024/1024),
            KB: Math.round(size/1024)
        }
    };

    return core
});