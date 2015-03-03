/**
 * Created by languid on 3/3/15.
 */
define([
'kernel',
'angular',
'./module'
],
function( core, ng, service){

    service.factory('GlobalControlUI', function() {
        var GlobalControlUI = core.Class.extend({
            init: function(){
                var self = this;

                this.element = $('#GlobalControl');
                this.$cancelBtn = this.element.find('.cancel');
                this.$nextBtn = this.element.find('.next');
                this.$content = this.element.find('p');

                this.nextClick = $.noop;
                this.cancelClick = $.noop;

                this.$nextBtn.click(function(e){
                    self.nextClick.call(self, e, this);
                });

                this.$cancelBtn.click(function(e){
                    self.cancelClick.call(self, e, this);
                });
            },
            show: function(text, next, cancel ){
                this.$content.html(text);
                if( $.isFunction(next) ){
                    this.nextClick = next;
                }
                if( $.isFunction(cancel) ){
                    this.cancelClick = cancel;
                }
                this._show();
            },
            _show: function(){
                this.element.addClass('show');
            },
            hide: function(){
                this.element.removeClass('show');
            }
        });

        return new GlobalControlUI;
    })
});