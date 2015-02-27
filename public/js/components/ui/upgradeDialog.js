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
        onclick: function( item, e ){
            $(e.currentTarget).toggleClass('checked');
        },
        render: function(){
            return (
                React.createElement("ul", null, 
                    this.props.list.map(function(item, index){
                        return (
                            React.createElement("li", {onClick: this.onclick.bind(this, item)}, 
                                React.createElement("span", {className: "action " + this.getAction(item.Action).toLowerCase()}, this.getAction(item.Action)), 
                                React.createElement("p", {className: "path"}, item.Path), 
                                React.createElement("span", {className: "checkbox"}, 
                                    React.createElement("i", {className: "fa fa-check"})
                                )

                            )
                        )
                    }, this)
                )
            )
        }
    });



    return function( events ){
        var buttons = [{
            text: 'Confirm',
            className: 'btn-primary',
            click: events.confirm
        }];

        var dialog = new Dialog({
            title: 'Upgrade',
            classStyle: 'upgrade-dialog'
        });
        dialog.upfileList = React.render(React.createElement(UpfileList, null), dialog.body[0]);

        dialog.formBtns = React.render(React.createElement(FormBtns, {
            buttons: buttons, 
            overload: dialog}
        ), dialog.footer[0]);

        return dialog;
    }
});