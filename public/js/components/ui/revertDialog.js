/**
 * Created by languid on 3/13/15.
 */
define([
'kernel',
'react',
'ui/Dialog'
],
function( core, React, Dialog ){
    var RevertList = React.createClass({displayName: "RevertList",
        getDefaultProps: function(){
            return {
                list: []
            }
        },
        more: function(){
            this.$list.find('li.hidden').removeClass('hidden')
        },
        render: function(){

            return (
                React.createElement("div", null, 
                    React.createElement("ul", null, 
                    this.props.list.map(function(item, index){
                        return (
                            React.createElement("li", {key: index, className:  index > 4 ? 'hidden' : ''}, 
                                React.createElement("span", {className: "path"}, item), 
                                React.createElement("span", {className: "control"}, 
                                    React.createElement("i", {className: "fa fa-download overwrite", onClick: this.props.events.revert.bind(this.props.overload, item)}), 
                                    React.createElement("i", {className: "fa fa-remove remove", onClick: this.props.events.remove.bind(this.props.overload, item)})
                                )
                            )
                        )
                    }, this)
                    ), 
                    React.createElement("p", {className:  'more ' + (this.props.list.length < 5 ? 'hidden': ''), onClick: this.more}, "Show More")
                )
            )
        },
        $el: null,
        getView: function(){
            this.$el = $(this.getDOMNode());
            this.$list = this.$el.find('ul');
        },
        setList: function( list ){
            this.setProps({
                list: list
            });
            this.forceUpdate();
            return this
        }
    });

    return function( events, options, extral ){
        events = $.extend({
            remove: $.noop,
            revert: $.noop
        }, events);

        var dialog = new Dialog( $.extend(options, {
            title: 'Revert Manager',
            classStyle: 'revert-dialog'
        }), extral);

        dialog.revertList = React.render(React.createElement(RevertList, {overload: dialog, events: events}), dialog.body[0], function(){
            this.getView();
        });

        return dialog;
    }
});