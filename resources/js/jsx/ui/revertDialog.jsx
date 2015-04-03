/**
 * Created by languid on 3/13/15.
 */
define([
'kernel',
'react',
'ui/Dialog'
],
function( core, React, Dialog ){
    var RevertList = React.createClass({
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
                <div>
                    <ul>
                    {this.props.list.map(function(item, index){
                        return (
                            <li key={index} className={ index > 4 ? 'hidden' : '' }>
                                <span className="path">{item}</span>
                                <span className="control">
                                    <i className="fa fa-download overwrite" onClick={this.props.events.revert.bind(this.props.overload, item)}></i>
                                    <i className="fa fa-remove remove" onClick={this.props.events.remove.bind(this.props.overload, item)}></i>
                                </span>
                            </li>
                        )
                    }, this)}
                    </ul>
                    <p className={ 'more ' + (this.props.list.length < 5 ? 'hidden': '') } onClick={this.more}>Show More</p>
                </div>
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

        dialog.revertList = React.render(<RevertList overload={dialog} events={events} />, dialog.body[0], function(){
            this.getView();
        });

        return dialog;
    }
});