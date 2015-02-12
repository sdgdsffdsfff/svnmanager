define(['kernel', 'react'], function( core, React ){
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
        onclick: function( item ){
            console.log(item)
        },
        render: function(){
            return (
                <div className="mod upfile-list">
                    <div className="hd">Upgrade</div>
                    <div className="bd">
                        <ul>
                            {this.props.list.map(function(item, index){
                                return (
                                    <li onClick={this.onclick.bind(this, item)}>
                                        <label className="checkbox">
                                            <input type="checkbox" />
                                            <span><i className="fa fa-check"></i></span>
                                        </label>
                                        <span className={"action " + this.getAction(item.Action).toLowerCase()}>{this.getAction(item.Action)}</span>
                                        <p className="path">{item.Path}</p>
                                    </li>
                                )
                            }, this)}
                        </ul>
                    </div>
                </div>
            )
        }
    });

    return function(){
        var div = $('<div />', {
            id: 'id'+core.random(10)
        });
        return div.data('reactElement', React.render(<UpfileList />, div[0]));
    };
});