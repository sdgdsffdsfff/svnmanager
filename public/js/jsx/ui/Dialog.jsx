define([
'kernel',
'react'
], function( core, React ){
    var Dialog = React.createClass({
        render: function(){
            return (
                <div className="modal-dialog">
                    <div className="modal-content">
                        <div className="modal-header">
                            <button type="button" className="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                            <h4 className="modal-title">{this.props.title}</h4>
                        </div>
                        <div className="modal-body"></div>
                        <div className="modal-footer"></div>
                    </div>
                </div>
            )
        }
    });

    return Dialog;
});