define([
'kernel',
'ui/Flyout'
],
function(core, Flyout){
    return function(anchor, html, options){
        anchor = $(anchor);

        var flyout = anchor.data('flyout');

        if( flyout == null ){
            var setting = {
                placement: 'top',
                alignment: 'center',
                destroy : true,
                stayTime : 800,
                classStyle: 'info',
                onHide: function(){
                    anchor.data('flyout', null);
                }
            };

            if( typeof options == 'string' && options.length ){
                setting.classStyle = 'alert-'+options;
            } else if( $.isPlainObject(options) || typeof options == 'undefined' ){
                $.extend(setting, options);
                if( setting.classStyle.length ){
                    setting.classStyle = 'alert-'+setting.classStyle;
                }
            }

            flyout = new Flyout('<div class="flyout-tips alert"></div>', setting);
            flyout.element.html(html).mouseenter(function() {
                anchor.data('flyout', flyout);
                flyout._clearStay()
            }).mouseleave(function() {
                anchor.data('flyout', null);
                flyout._createStayTimer()
            });
            flyout.arrow();
            flyout.show(anchor);
        }

        return flyout;
    }
});