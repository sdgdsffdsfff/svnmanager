/**
 * Created by languid on 3/12/15.
 */

define(['kernel'], function( core ){
    var Spinner = core.Class.extend({
        init: function(element, options){
            this.$el = element;
            this.options = $.extend({}, Spinner.DefaultOptions, options)
            this.isLoading = false;
        },
        setState: function( opt ){
            var $el = this.$el;
            var data = $el.data();

            if( data.resetClass == null ){
                $el.data('resetClass', $el.attr('class'))
            }

            if( opt == 'loading' ){
                this.isLoading = true;
                $el.attr('class', this.options.loadingClass)
            } else {
                if( this.isLoading ){
                    this.isLoading = false;
                    $el.attr('class', data.resetClass)
                }
            }
        }
    }, {
        DefaultOptions: {
            loadingClass: 'fa fa-spinner fx-spinner'
        }
    });

    function Plugin(option) {
        return this.each(function () {
            var $this   = $(this);
            var data    = $this.data('ui.spinner');
            var options = typeof option == 'object' && option;

            if (!data) $this.data('ui.spinner', (data = new Spinner($this, options)));

            if (option == 'toggle') data.toggle();
            else if (option) data.setState(option)
        })
    }

    var old = $.fn.spinner;

    $.fn.spinner = Plugin;
    $.fn.spinner.Constructor = Spinner;

    $.fn.spinner.noConflict = function () {
        $.fn.spinner = old;
        return this
    };
});