define([
'kernel',
'react',
'ui/Dialog',
'components/form/FormBtns'
],
function( core, React, Dialog, FormBtns ){
    var UpfileList = React.createClass({
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
                <ul>
                    {this.props.list.map(function(item, index){
                        return (
                            <li onClick={this.onclick.bind(this, item)}>
                                <span className={"action " + this.getAction(item.Action).toLowerCase()}>{this.getAction(item.Action)}</span>
                                <p className="path">{item.Path}</p>
                                <span className="checkbox">
                                    <i className="fa fa-check"></i>
                                </span>

                            </li>
                        )
                    }, this)}
                </ul>
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
        dialog.upfileList = React.render(<UpfileList />, dialog.body[0]);

        dialog.formBtns = React.render(<FormBtns
            buttons={buttons}
            overload={dialog}
        />, dialog.footer[0]);

        return dialog;
    }
});