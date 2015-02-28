define([
'kernel',
'ui/Flyout',
'components/ui/confirm'
],
function(core, Flyout, confirm){
    var flyout = null, text = null;

    return function(anchor, ok, options){

        if( typeof options == 'string' ){
            options = {
                html: options
            }
        }

        options = $.extend({
            okText: 'Confirm',
            okClass: 'btn-primary',
            cancelText: 'Cancel',
            cancelClass: 'btn-default',
            html: 'Are you sure?',
            okClick: ok || $.noop
        }, options);

        var elem = confirm(options);

        if(flyout == null){
            flyout = new Flyout(elem, {
                onRendered: function(element){
                    var thisFlyout = this;
                    element.delegate('button[role]','click',function(){
                        thisFlyout.hide();
                        var type = $(this).attr('role'), anchor = thisFlyout._currentAnchor;
                        if(type == 'ok'){
                            options.okClick.call( thisFlyout, anchor )
                        }
                        return false
                    });
                }
            });
            text = flyout.element.find('p.text');

            core.body.append( elem );
        }

        if( options.html ){
            text.css('width','auto').show().html( options.html );
            if(flyout.element.outerWidth() > 300 ){
                text.width(250);
            }
        }else{
            text.hide().html('');
        }

        flyout.show( anchor, "top");

        return flyout
    }
});