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
                onHide: function(){
                    anchor.data('flyout', null);
                }
            };

            if( $.isPlainObject(options) ){
                $.extend(setting, options);
            }

            flyout = new Flyout('<div class="info"></div>', setting);

            flyout.element.html(html).mouseenter(function() {
                anchor.data('flyout', flyout);
                flyout._clearStay()
            }).mouseleave(function() {
                anchor.data('flyout', null);
                flyout._createStayTimer()
            });

            flyout.show(anchor);
        }

        return flyout;
    }
});