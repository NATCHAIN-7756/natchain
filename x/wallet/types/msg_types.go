package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgRegisterWallet 注册矿机消息
type MsgRegisterWallet struct {
	Creator string `json:"creator"`
	CpuSn   string `json:"cpu_sn"`
}

func NewMsgRegisterWallet(creator, cpuSn string) *MsgRegisterWallet {
	return &MsgRegisterWallet{Creator: creator, CpuSn: cpuSn}
}

func (m *MsgRegisterWallet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid creator: %s", err)
	}
	if m.CpuSn == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("cpu_sn cannot be empty")
	}
	return nil
}

func (m *MsgRegisterWallet) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Creator)
	return []sdk.AccAddress{addr}
}

type MsgRegisterWalletResponse struct {
	Index string `json:"index"`
}

// MsgHeartbeat 心跳消息
type MsgHeartbeat struct {
	Creator string `json:"creator"`
	Index   string `json:"index"`
}

func NewMsgHeartbeat(creator, index string) *MsgHeartbeat {
	return &MsgHeartbeat{Creator: creator, Index: index}
}

func (m *MsgHeartbeat) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid creator: %s", err)
	}
	if m.Index == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("index cannot be empty")
	}
	return nil
}

func (m *MsgHeartbeat) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Creator)
	return []sdk.AccAddress{addr}
}

type MsgHeartbeatResponse struct {
	Success bool `json:"success"`
}
