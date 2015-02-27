define([
'kernel',
'react',
'ui/Dialog',
'components/form/FormBtns'
],
function( core, React, Dialog, FormBtns ){
    var UpfileList = React.createClass({displayName: "UpfileList",
        getDefaultProps: function(){
            return {
                list: []
            }
        },
        setList: function( list ){
            this.setProps({
                list: list
            });
            this.forceUpdate();
        },
        getAction: function( type ){
            return type == 1 ? 'Add' :
                   type == 2 ? 'Update' :
                   type == 3 ? 'Delete' : '';
        },
        selectAll: function( e ){
            $(this.getDOMNode()).find('input:checkbox').attr('checked', 'checked');
            $(e.target).html('Deselect All');
        },
        render: function(){
            return (
                React.createElement("div", null, 
                    React.createElement("ul", null, 
                        this.props.list.map(function(item, index){
                            return (
                                React.createElement("li", {key: index}, 
                                    React.createElement("label", null, 
                                        React.createElement("span", {className: "action " + this.getAction(item.Action).toLowerCase()}, this.getAction(item.Action)), 
                                        React.createElement("span", {className: "path"}, item.Path), 
                                        React.createElement("input", {type: "checkbox", className: "hidden"}), 
                                        React.createElement("span", {className: "checkbox"}, 
                                            React.createElement("i", {className: "fa fa-check"})
                                        )
                                    )
                                )
                            )
                        }, this)
                    ), 
                    React.createElement("p", {className: "control"}, 
                        React.createElement("span", null, "Tree View"), 
                        React.createElement("span", {onClick: this.selectAll}, "Select All")
                    )
                )
            )
        }
    });

    return function( events, options, extral ){
        var buttons = [{
            text: 'Deploy Now',
            className: 'btn-primary',
            click: events.confirm
        }];

        var dialog = new Dialog( $.extend(options, {
            title: 'Undeploy Application',
            classStyle: 'upgrade-dialog'
        }), extral);
        dialog.upfileList = React.render(React.createElement(UpfileList, null), dialog.body[0]);

        dialog.formBtns = React.render(React.createElement(FormBtns, {
            buttons: buttons, 
            overload: dialog}
        ), dialog.footer[0]);

        return dialog;
    }
});