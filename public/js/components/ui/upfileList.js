define([
'kernel',
'react',
'ui/Dialog',
'components/FormBtns'
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
        onclick: function( item ){
            console.log(item)
        },
        render: function(){
            return (
                React.createElement("ul", null, 
                    this.props.list.map(function(item, index){
                        return (
                            React.createElement("li", {onClick: this.onclick.bind(this, item)}, 
                                React.createElement("label", {className: "checkbox"}, 
                                    React.createElement("input", {type: "checkbox"}), 
                                    React.createElement("span", null, React.createElement("i", {className: "fa fa-check"}))
                                ), 
                                React.createElement("span", {className: "action " + this.getAction(item.Action).toLowerCase()}, this.getAction(item.Action)), 
                                React.createElement("p", {className: "path"}, item.Path)
                            )
                        )
                    }, this)
                )
            )
        }
    });

    var buttons = [{
        text: 'Confirm',
        className: 'btn-primary',
        click: function () {
            this.submitForm();
        }
    }];

    var dialog = new Dialog();
    dialog.upfileList = React.render(React.createElement(UpfileList, null), dialog.body[0]);

    dialog.formBtns = React.render(React.createElement(FormBtns, {
        buttons: buttons, 
        overload: dialog}
    ), dialog.footer[0]);

    return dialog;
});