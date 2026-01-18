"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Modal } from "@/components/ui/modal";
import { useCreateAccount } from "@/hooks/useAccounts";
import React, { useState } from "react";

interface AddAccountModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export const AddAccountModal = ({ isOpen, onClose, onSuccess }: AddAccountModalProps) => {
  const [formData, setFormData] = useState({
    name: "",
    type: "depository",
    subtype: "checking",
    balance: "",
    currency: "USD",
  });
  const { mutate: createAccount, isPending: loading, error: mutationError } = useCreateAccount();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    createAccount(
      {
        ...formData,
        balance: parseFloat(formData.balance),
      },
      {
        onSuccess: () => {
          onSuccess();
          onClose();
          setFormData({
            name: "",
            type: "depository",
            subtype: "checking",
            balance: "",
            currency: "USD",
          });
        },
      }
    );
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Connect Manual Account">
      <form onSubmit={handleSubmit} className="space-y-6">
        <Input
          label="Account Name"
          placeholder="e.g. Chase Checking"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          required
        />

        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-1.5">
            <label className="block text-sm font-medium text-zinc-400 ml-1" htmlFor="account-type-select">Type</label>
            <select
              id="account-type-select"
              title="Select Account Type"
              className="w-full bg-brand-muted border border-zinc-800 rounded-2xl px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-brand-lime/20 transition-all appearance-none"
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value })}
              required
            >
              <option value="depository">Cash / Bank</option>
              <option value="credit">Credit Card</option>
              <option value="investment">Investment</option>
              <option value="loan">Loan / Mortgage</option>
            </select>
          </div>
          <Input
            label="Initial Balance ($)"
            type="number"
            step="0.01"
            placeholder="0.00"
            value={formData.balance}
            onChange={(e) => setFormData({ ...formData, balance: e.target.value })}
            required
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Input
            label="Subtype"
            placeholder="e.g. checking, savings"
            value={formData.subtype}
            onChange={(e) => setFormData({ ...formData, subtype: e.target.value })}
          />
          <Input
            label="Currency"
            placeholder="USD"
            value={formData.currency}
            onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
          />
        </div>

        {mutationError && (
          <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-red-500 text-sm">
            {(mutationError as Error).message}
          </div>
        )}

        <div className="space-y-4 pt-4">
          <Button type="submit" className="w-full py-4" disabled={loading}>
            {loading ? "Adding Account..." : "Add Account"}
          </Button>
          <Button variant="ghost" className="w-full py-4 text-zinc-500" onClick={onClose} type="button">
            Cancel
          </Button>
        </div>
      </form>
    </Modal>
  );
};
