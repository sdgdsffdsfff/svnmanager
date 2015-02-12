define([
'kernel',
'react'
], function( core, React ){
    var Confirm = React.createClass({
        render: function(){
            return (
                <div>
                    <p className="text"></p>
                    <div className="buttons">
                        <button className={"btn btn-xs "+this.props.okClass} role="ok">{this.props.okText}</button>
                        <button className={"btn btn-xs "+this.props.cancelClass} role="cancel">{this.props.cancelText}</button>
                    </div>
                </div>
            )
        }
    });

    return Confirm;
});