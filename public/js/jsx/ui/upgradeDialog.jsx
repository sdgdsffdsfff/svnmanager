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
        selectAll: function( e ){
            $(this.getDOMNode()).find('input:checkbox').attr('checked', 'checked');
            $(e.target).html('Deselect All');
        },
        render: function(){
            return (
                <div>
                    <ul>
                        {this.props.list.map(function(item, index){
                            return (
                                <li key={index}>
                                    <label>
                                        <span className={"action " + this.getAction(item.Action).toLowerCase()}>{this.getAction(item.Action)}</span>
                                        <span className="path">{item.Path}</span>
                                        <input type="checkbox" className="hidden" />
                                        <span className="checkbox">
                                            <i className="fa fa-check"></i>
                                        </span>
                                    </label>
                                </li>
                            )
                        }, this)}
                    </ul>
                    <p className="control">
                        <span>Tree View</span>
                        <span onClick={this.selectAll}>Select All</span>
                    </p>
                </div>
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
        dialog.upfileList = React.render(<UpfileList />, dialog.body[0]);

        dialog.formBtns = React.render(<FormBtns
            buttons={buttons}
            overload={dialog}
        />, dialog.footer[0]);

        return dialog;
    }
});